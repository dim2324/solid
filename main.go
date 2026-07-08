package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Order представляет заказ
type Order struct {
	ID       int
	Customer string
	Products []string
	Total    float64
	Status   string
}

// RepositoryInitializer - интерфейс для инициализации базы данных (ISP)
type RepositoryInitializer interface {
	Initialize() error
}

// RepositoryWriter - интерфейс для записи данных (ISP)
type RepositoryWriter interface {
	SaveOrder(order Order) error
}

// Notifier - базовый интерфейс для отправки уведомлений (LSP)
type Notifier interface {
	Send(customer string, message string) error
}

// SQLiteOrderRepository реализует оба интерфейса для работы с БД
type SQLiteOrderRepository struct {
	db *sql.DB
}

func NewSQLiteOrderRepository(db *sql.DB) *SQLiteOrderRepository {
	return &SQLiteOrderRepository{db: db}
}

// Initialize реализует RepositoryInitializer
func (r *SQLiteOrderRepository) Initialize() error {
	_, err := r.db.Exec(`
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer TEXT NOT NULL,
        products TEXT NOT NULL,
        total REAL NOT NULL,
        status TEXT NOT NULL
    )`)
	return err
}

// SaveOrder реализует RepositoryWriter
func (r *SQLiteOrderRepository) SaveOrder(order Order) error {
	productsStr := strings.Join(order.Products, ", ")
	_, err := r.db.Exec(
		"INSERT INTO orders (customer, products, total, status) VALUES (?, ?, ?, ?)",
		order.Customer, productsStr, order.Total, order.Status,
	)
	return err
}

// EmailSender реализует интерфейс Notifier
type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

// Send реализует метод интерфейса Notifier
func (e *EmailSender) Send(customer string, message string) error {
	fmt.Printf("Email уведомление отправлено клиенту %s: %s\n", customer, message)
	return nil
}

// SMSSender реализует интерфейс Notifier
type SMSSender struct{}

func NewSMSSender() *SMSSender {
	return &SMSSender{}
}

// Send реализует метод интерфейса Notifier
func (s *SMSSender) Send(customer string, message string) error {
	fmt.Printf("SMS уведомление отправлено клиенту %s: %s\n", customer, message)
	return nil
}

// OrderService - основной сервис, зависящий от абстракций (DIP)
type OrderService struct {
	repo     RepositoryWriter
	notifier Notifier
}

func NewOrderService(repo RepositoryWriter, notifier Notifier) *OrderService {
	return &OrderService{
		repo:     repo,
		notifier: notifier,
	}
}

// CreateOrder - бизнес-логика создания заказа (SRP)
func (s *OrderService) CreateOrder(customer string, products []string, total float64) error {
	order := Order{
		Customer: customer,
		Products: products,
		Total:    total,
		Status:   "pending",
	}

	// Сохраняем заказ в БД
	if err := s.repo.SaveOrder(order); err != nil {
		return fmt.Errorf("ошибка сохранения заказа: %w", err)
	}

	// Отправляем уведомление
	message := fmt.Sprintf("Ваш заказ на сумму %.2f создан и находится в обработке", total)
	if err := s.notifier.Send(customer, message); err != nil {
		return fmt.Errorf("ошибка отправки уведомления: %w", err)
	}

	return nil
}

func main() {
	// Инициализация базы данных
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем репозиторий
	repo := NewSQLiteOrderRepository(db)

	// Инициализируем таблицы
	if err := repo.Initialize(); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	// Демонстрация с Email уведомлениями
	fmt.Println("=== Email ===")
	emailService := NewOrderService(repo, NewEmailSender())
	err = emailService.CreateOrder("Иван", []string{"apple", "banana"}, 10.5)
	if err != nil {
		log.Fatal("Ошибка создания заказа:", err)
	}

	// Демонстрация с SMS уведомлениями (OCP - расширение без изменения кода)
	fmt.Println("\n=== SMS ===")
	smsService := NewOrderService(repo, NewSMSSender())
	err = smsService.CreateOrder("Мария", []string{"orange", "grape"}, 15.75)
	if err != nil {
		log.Fatal("Ошибка создания заказа:", err)
	}

	// Демонстрация с множественными уведомлениями
	fmt.Println("\n=== Email + SMS ===")
	multiService := NewOrderService(repo, NewMultiNotifier(
		NewEmailSender(),
		NewSMSSender(),
	))
	err = multiService.CreateOrder("Петр", []string{"milk", "bread"}, 5.25)
	if err != nil {
		log.Fatal("Ошибка создания заказа:", err)
	}
}

// MultiNotifier позволяет отправлять уведомления через несколько каналов (OCP)
type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

func (m *MultiNotifier) Send(customer string, message string) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Send(customer, message); err != nil {
			return err
		}
	}
	return nil
}
