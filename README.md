# Subserv

**Subserv** is a lightweight subscription service backend built in **Go** using the **Gin** web framework. It provides a foundation for managing user subscriptions, including authentication, plan management, and payment integration.

## 🚀 Features

- 📦 Subscription plan management
- 🧪 Unit tests
- 🛠️ Modular and extensible architecture
- 📄 OpenAPI documentation
- 📊 Swagger UI for API exploration
- 🔐 JWT-based authentication

## 📦 Tech Stack

- **Go** (Golang)
- **Gin** – HTTP web framework
- **GORM** – ORM for database interactions
- **Sqlite** – Default database
- **Swagger** – API documentation and testing
- **Testify** – Testing framework
## 🛠️ Project Description
The development is based on the following user story:
As a User, I want to be able to select from a list a product, and based on this product to receive a subscription plan:
AC:
1. I can fetch a list of products
2. I can fetch a single product
3. I can buy a single product
4. I want to fetch the following informations related to my subscription (e.g. start date, end date, duration of the subscription, prices, tax)
5. I can pause and unpause my subscription
6. I can cancel my active subscription

## 📝 Remarks

- Only the endpoints related to the user story are implemented. (e.g. no admin endpoints, user management, payment management, etc.)
- The application is designed to be modular and extensible, allowing for easy addition of new features and endpoints in the future.
- **Important** The application uses JWT for authentication, but does not implement user management or registration endpoints, so in the swagger UI use the `Authorization` header to pass the JWT token for testing purposes. Simply pass **`Bearer test-token`** as the value of the `Authorization` header in your requests.
- Payment is a dummy implementation and does not involve real payment processing. The payment processor is designed to simulate a successful payment transaction for testing purposes with %5 chance of failure.
- Unit tests are provided to ensure the functionality of the application. The tests cover the main features and endpoints, but do not include exhaustive coverage of all possible scenarios. and integration tests are not implemented.
- Current pause logic is a simple implementation that does not account for complex scenarios such as overlapping pause periods or multiple pauses. It is designed to demonstrate the basic functionality of pausing and resuming subscriptions.
- Docker, Makefile, and other common development tools are not used in this project to keep the implementation simple and focused on the core functionality. However, the project can be easily extended to include these tools in the future if needed.
- Configurations are hardcoded in the codebase for simplicity, but usually they are implemented by Viper and managed by environment variables or configuration files in production applications.
- A proper logging implementation is not included in this project. The application uses simple print statements for logging, but in a production application, a structured logging library (my choice being **Logrus**) would be used to provide better logging capabilities.
## 📦 Setup
If you want to run it through Docker:
```bash
docker run -it -p 8080:8080 golang:1.24.4-bookworm /bin/sh
```
Fetch repository, and run the following command to install the dependencies:
```bash
git clone https://github.com/thatmatin/subserv.git
cd subserv
go mod tidy
```

## 🏃‍♂️ Running the application
In order to populate the database with initial data, run the following command:

```bash
go run . populate
```
To explore the API documentation, you can visit the Swagger UI To launch the Swagger UI, run the following command in your terminal:

```bash
go run . serve -s
```
Then, open your web browser and navigate to `http://localhost:8080/swagger/index.html` to view the API documentation and test the endpoints interactively.
This command will start the server and serve the Swagger UI at the specified address. (You can drop the `-s` flag if you want to run the server without serving Swagger UI.)

## 🧪 Running tests
To run the tests, use the following command:

```bash
go test ./... -v
```

## To Alejandro
Best of luck. I hope you find this project useful and inspiring. If you have any questions or suggestions, you have my contact ;)
