# Документация REST API для получения остатков товара 

### Базовый URL

`https://covrytrade.ru/exchange/api/v1`

## Точки доступа (Endpoints)

### 1. Коллекции

### 1.1 **Cписок коллекций**

**Метод:** `GET`

**URL:** `/collection`

**Параметры запроса:**

| Parameter     | Type     | Required | Description                        |
|---------------|----------|----------|------------------------------------|
| `collection`  | `string` |          | Коллекция                          |
| `uuid`        | `string` |          | Уникальный идентификатор коллекции |
| `code`        | `string` |          | Код коллекции                      |
| `description` | `string` |          | Название коллекции                 |
| `length`      | `string` |          | Длина                              |
| `width`       | `string` |          | Ширина                             |

**Пример ответа:**

```json
[
  {
    "collection": "Тафтинговые покрытия и дорожка 5 цветные",
    "uuid": "4c8d080027792da011e8a6177df7379c",
    "code": "2685",
    "description": "Тафтинговые покрытия и дорожка 5 цветные (0,7x21,5)",
    "length": "21.50",
    "width": "0.70",
    "count": 1
  },
  {
    "collection": "Тафтинговые покрытия и дорожка 5 цветные",
    "uuid": "e597080027792da011eae5d2db16ee02",
    "code": "8370",
    "description": "Тафтинговые покрытия и дорожка 5 цветные (1,3x18,5)",
    "length": "18.50",
    "width": "1.30",
    "count": 1
  }
]
```

**Параметры ответа JSON:**

| Parameter     | Type     | Required | Description                        |
|---------------|----------|----------|------------------------------------|
| `collection`  | `string` |          | Коллекция                          |
| `uuid`        | `string` |          | Уникальный идентификатор коллекции |
| `code`        | `string` |          | Код коллекции                      |
| `description` | `string` |          | Название коллекции                 |
| `length`      | `string` |          | Длина                              |
| `width`       | `string` |          | Ширина                             |
| `count`       | `number` |          | Доступный остаток на складе        |

---

### 2. Характеристики

### 2.1 **Список характеристик**

**Метод:** `GET`

**URL:** `/details`

**Параметры запроса:**

| Parameter     | Type     | Required | Description                             |
|---------------|----------|----------|-----------------------------------------|
| `collection`  | `string` |          | Коллекция                               |
| `uuid`        | `string` |          | Уникальный идентификатор коллекции      |
| `code`        | `string` |          | Код коллекции                           |
| `description` | `string` |          | Название коллекции                      |
| `length`      | `string` |          | Длина                                   |
| `width`       | `string` |          | Ширина                                  |
| `uuidDetails` | `string` |          | Уникальный идентификатор характеристики |
| `codeDetails` | `string` |          | Код характеристики                      |
| `picture`     | `string` |          | Характеристика изображения              |
| `form`        | `string` |          | Форма                                   |
| `color`       | `string` |          | Цвет                                    |
| `brand`       | `string` |          | Производитель                           |
| `barcode`     | `string` |          | Штрихкод                                |

**Пример ответа:**

```json
[
  {
    "collection": "Viva",
    "uuid": "4597080027792da011e8898f6887266e",
    "code": "193",
    "description": "Viva (0,7x1,4)",
    "length": "1.40",
    "width": "0.70",
    "count": 9,
    "uuidDetails": "188b080027792da011e9deaa84624e4c",
    "codeDetails": "163877",
    "picture": "1039",
    "form": "1",
    "color": "32100",
    "brand": "Moldabela",
    "barcode": "4841227940184",
    "countBalance": 10,
    "reservedBalance": 1
  },
  {
    "collection": "Viva",
    "uuid": "4597080027792da011e8898f6887266e",
    "code": "193",
    "description": "Viva (0,7x1,4)",
    "length": "1.40",
    "width": "0.70",
    "count": 2,
    "uuidDetails": "4284080027792da011e9335190a34de6",
    "codeDetails": "105997",
    "picture": "1039",
    "form": "1",
    "color": "34700",
    "brand": "Moldabela",
    "barcode": "4841808020908",
    "countBalance": 2,
    "reservedBalance": 0
  }
]
```

**Параметры ответа JSON:**

| Parameter         | Type     | Required | Description                             |
|-------------------|----------|----------|-----------------------------------------|
| `collection`      | `string` |          | Коллекция                               |
| `uuid`            | `string` |          | Уникальный идентификатор коллекции      |
| `code`            | `string` |          | Код коллекции                           |
| `description`     | `string` |          | Название коллекции                      |
| `length`          | `string` |          | Длина                                   |
| `width`           | `string` |          | Ширина                                  |
| `uuidDetails`     | `string` |          | Уникальный идентификатор характеристики |
| `codeDetails`     | `string` |          | Код характеристики                      |
| `picture`         | `string` |          | Характеристика изображения              |
| `form`            | `string` |          | Форма                                   |
| `color`           | `string` |          | Цвет                                    |
| `brand`           | `string` |          | Производитель                           |
| `barcode`         | `string` |          | Штрихкод                                |
| `countBalance`    | `string` |          | Фактический остаток на складе           |
| `reservedBalance` | `string` |          | Количество зарезервированного товара    |
| `count`           | `string` |          | Доступный остаток (Факт - Резерв)       |

---

### 3. Изображения

### 3.1 **Images**

**Метод:** `GET`

**URL:** `/images`

**Параметры запроса:**

| Parameter    | Type     | Required | Description                |
|--------------|----------|----------|----------------------------|
| `collection` | `string` |          | Коллекция                  |
| `picture`    | `string` |          | Характеристика изображения |
| `form`       | `string` |          | Форма                      |
| `color`      | `string` |          | Цвет                       |
| `brand`      | `string` |          | Производитель              |

**Пример ответа:**

```json
[
  {
    "logoURL": "https://covrytrade.ru/exchange/files/e86fab123f1f11e9aa8c080027792da0/logo/e86fab123f1f11e9aa8c080027792da0.png",
    "imagesURL": null,
    "collection": "MALAGA",
    "picture": "12001",
    "form": "1",
    "color": "120",
    "brand": "Carpetoff"
  },
  {
    "logoURL": "https://covrytrade.ru/exchange/files/057dad943f2011e9aa8c080027792da0/logo/057dad943f2011e9aa8c080027792da0.png",
    "imagesURL": null,
    "collection": "MALAGA",
    "picture": "12003",
    "form": "1",
    "color": "120",
    "brand": "Carpetoff"
  }
]
```

**Параметры ответа JSON:**

| Parameter    | Type     | Required | Description                  |
|--------------|----------|----------|------------------------------|
| `logoURL`    | `string` |          | ссылка на изображение, url   |
| `imagesURL`  | `string` |          | Массив изображений URL(план) |
| `collection` | `string` |          | Коллекция                    |
| `picture`    | `string` |          | Характеристика изображения   |
| `form`       | `string` |          | Форма                        |
| `color`      | `string` |          | Цвет                         |
| `brand`      | `string` |          | Производитель                |
