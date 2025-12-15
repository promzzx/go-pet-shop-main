# Postman-набор для Go Pet Shop

В каталоге лежит коллекция и окружение для быстрой проверки всех HTTP-хендлеров проекта.

## Файлы
- `go-pet-shop.postman_collection.json` — коллекция запросов с автотестами на статус/JSON.
- `go-pet-shop.postman_environment.json` — окружение с базовым URL и тестовыми ID.

## Как использовать
1. Запустите приложение (`go run cmd/app/main.go`) и убедитесь, что оно доступно на `http://localhost:3000` (или поправьте переменную `base_url` в окружении).
2. Импортируйте оба файла в Postman (File → Import).
3. Выберите окружение **Go Pet Shop Local**.
4. Выполняйте запросы по разделам: Health, Products, Customers, Orders, Checkout, History. В теле запросов уже проставлены примерные JSON-данные.
5. При необходимости смените значения переменных `productId`, `secondProductId`, `customerId`, `customerEmail`, `orderId` в окружении под ваши тестовые данные.

## Как проверить хендлеры в Postman
Ниже — короткий сценарий проверки прямо в Postman, чтобы увидеть, что автотесты коллекции проходят:

1. **Health** → `GET /health`: нажмите *Send* и убедитесь, что тесты внизу вкладки показывают зелёный статус (`status code is 200`).
2. **Products**:
   - `POST /products`: выполните запрос с примерным телом, скопируйте `id` из ответа и подставьте его в переменную `productId` окружения.
   - `GET /products/:id`: повторите запрос с обновлённым `productId` и убедитесь, что тесты «status code is 200» и «response has product id» проходят.
3. **Customers**: выполните `POST /customers` (при необходимости поправьте email), сохраните `id` в переменную `customerId`, затем вызовите `GET /customers/:id` и проверьте зелёные тесты.
4. **Orders**:
   - `POST /orders`: используйте `customerId` и `productId` в теле. Сохраните `id` ответа в переменную `orderId`.
   - `GET /orders/:id`: убедитесь, что оба теста проходят.
5. **Checkout**: вызовите `POST /checkout` с `orderId` и любым `paymentToken` — тест проверит `200`.
6. **History**: выполните `GET /history/:customerId` — тест проверит успешный статус.

### Запуск всех тестов сразу
Вы можете прогнать все запросы подряд через **Collection Runner**:
1. В дереве коллекции нажмите правой кнопкой `Go Pet Shop` → *Run collection*.
2. Выберите окружение **Go Pet Shop Local** и нажмите *Run*. Если заранее подставлены корректные переменные `productId`, `customerId`, `orderId`, вы увидите общее число пройденных тестов.