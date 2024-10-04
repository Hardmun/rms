# RMS - Remains Stocks API

## Overview

RMS (Remains Stocks) is a REST API service that provides real-time product balance information for a company. The API
offers various endpoints to access product collections, detailed balances, and image URLs based on specific criteria.

---

## Features

1. **Client Access Request**:
    - To access the API, clients must submit a form with mandatory fields: Company, E-mail, and Initials.
    - To prevent spam, the API server enforces rate limiting on requests.

2. **Email Confirmation**:
    - Once the form is submitted, an email is sent to the company containing a URL for access approval. This URL can be
      used only once and grants the client access credentials (username and password).

3. **Authentication**:
    - The service uses **Basic Authentication**. Client passwords are securely stored using one-way encryption in the
      local database (e.g., SQLite).
    - Once the user is authenticated, the session must store user-specific settings, but
      passwords should never be stored in the session.
    - User sessions must expire after a specified period of inactivity
      to ensure security and optimize memory consumption.

4. **API Resources**:
    - `/collection`: Returns a summary of product balances for specific collections.
        - Fields: `collection`, `uuid`, `code`, `description`, `length`, `width`, `count`.

    - `/details`: Provides detailed product balance information.
        - Fields: Collection fields plus additional info such as `UuidDetails`, `CodeDetails`, `Picture`, `Form`,
          `color`, `brand`, `barcode`, `countBalance`, `reservedBalance`.

    - `/images`: Returns product image URLs based on specific parameters (collection, picture, form, color, brand).

5. **Database**:
    - All data is retrieved from a **PostgreSQL** database. Database schema and table definitions should be described in
      a settings file.

6. **User Session Management**:
    - User sessions are restricted by specific configuration parameters, limiting the number of simultaneous connections
      a user can have to the server to control usage and prevent abuse.

7. **SSL Support**:
    - The service supports SSL/TLS if `.crt` and `.key` files are present in a designated folder. Otherwise, the
      connection will fall back to non-SSL.

---

## Requirements

- **PostgreSQL Database**: The service fetches data from a PostgreSQL database.
- **Basic Authentication**: Client credentials are securely stored with one-way encryption.
- **Linux (Ubuntu 22.04)**: The service must be deployed on a Linux-based system.

---

## Installation Instructions

1. **Prerequisites**:
    - Ensure that PostgreSQL is installed and running.
    - Install necessary dependencies (e.g., Go, Postgres driver).
    - Linux (Ubuntu 22.04) is required for deployment.

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/Hardmun/rms.git
   ```

3. **Configuration**:
    - Populate the `settings` file with the required parameters, such as database connection details, session limits,
      and other API configurations.
    - Ensure your `.crt` and `.key` files for SSL are placed in the designated folder.

4. **Database Setup**:
    - Create a PostgreSQL database and tables as per the provided schema in the settings file.

5. **Follow the instructions using the install.md file to deploy the service.**

---

## Logging

- All errors and critical information are logged in a specified log file.

---

## API Endpoints

### **/collection**

Returns product balance information for a specific collection.

**Response Example**:

```json
{
  "collection": "Carpets",
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "code": "COL123",
  "description": "Classic Carpet Collection",
  "length": "5m",
  "width": "3m",
  "count": 50.0
}
```

### **/details**

Returns detailed product information, including balance and additional metadata.

**Response Example**:

```json
{
  "collection": "Carpets",
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "code": "COL123",
  "description": "Classic Carpet Collection",
  "length": "5m",
  "width": "3m",
  "count": 50.0,
  "uuidDetails": "456e1234-e89b-12d3-a456-426614174001",
  "codeDetails": "DET123",
  "picture": "Geometric Pattern",
  "form": "Rectangular",
  "color": "Blue",
  "brand": "CarpetCo",
  "barcode": "1234567890123",
  "countBalance": 30.0,
  "reservedBalance": 5.0
}
```

### **/images**

Returns product image URL based on specific parameters (Collection, Picture, Form, Color, Brand).

**Response Example**:

```json
{
  "logoURL": "https://example.com/files/605028eaaf4411e7f189080027792da0/logo/605028eaaf4411e7f189080027792da0.png"
}
```

---

## Security Considerations

- Passwords are stored using secure, one-way encryption.
- API access is controlled through Basic Authentication.
- Optional SSL connection is supported.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

```