Примененные принципы SOLID:
1. Single Responsibility Principle (SRP)
SQLiteOrderRepository отвечает только за работу с БД

OrderService отвечает только за бизнес-логику

Каждый отправитель (EmailSender, SMSSender) отвечает только за свой тип уведомлений

2. Open/Closed Principle (OCP)
Система открыта для расширения (можно добавить новые типы БД и уведомлений)

Закрыта для модификации (не нужно менять существующий код для добавления нового функционала)

MultiNotifier демонстрирует расширение без изменения существующих классов

3. Liskov Substitution Principle (LSP)
Все реализации Notifier могут быть взаимозаменяемы

EmailSender, SMSSender, MultiNotifier могут использоваться везде, где ожидается Notifier

4. Interface Segregation Principle (ISP)
RepositoryInitializer и RepositoryWriter разделены на отдельные интерфейсы

Клиенты зависят только от тех методов, которые им нужны

5. Dependency Inversion Principle (DIP)
OrderService зависит от абстракций (RepositoryWriter, Notifier), а не от конкретных реализаций

Зависимости внедряются через конструктор

Преимущества новой архитектуры:
Расширяемость: легко добавить новые типы БД (PostgreSQL, MongoDB) или уведомлений (Push, Telegram)

Тестируемость: можно легко замокать зависимости для unit-тестов

Гибкость: можно комбинировать разные реализации в runtime

Поддерживаемость: каждый компонент имеет четкую ответственность
