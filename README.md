# ActualBudgetNormalizer

ActualBudgetNormalizer is a Go application that processes financial transaction data from a CSV file, augments it with additional information using an AI model, and saves the augmented data back to a CSV file.

## Features

- Loads transaction data from a CSV file
- Processes each transaction concurrently using goroutines
- Queries an AI model to extract relevant information from transaction descriptions
- Augments the transaction data with categories, creditors, and tags
- Saves the augmented data to a new CSV file
- Provides a loading animation during the processing

## Prerequisites

- Go 1.20 or higher
- Required dependencies (listed in `go.mod`)

## Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/your-username/ActualBudgetNormalizer.git
   ```

2. Change to the project directory:

   ```shell
   cd ActualBudgetNormalizer
   ```

3. Install the dependencies:

   ```shell
   go mod download
   ```

## Usage

1. Place your transaction data CSV file in the `testdata` directory.

2. Run the application:

   ```shell
   make run
   ```

   This will build and run the application, processing the transaction data and generating the augmented data CSV file.

3. The augmented data will be saved as `augmented_data.csv` in the `testdata` directory.

## Testing

To run the tests for the project, use the following command:

```shell
make test
```

This will run all the tests in the project and display the test results.

## Linting

To run the linter and check for code style and potential issues, use the following command:

```shell
make lint
```

This will run the `golangci-lint` linter and display any linting errors or warnings.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

```

This `README.md` file provides an overview of the project, its features, prerequisites, installation instructions, usage guidelines, testing and linting commands, and information about contributing and licensing.

Feel free to customize the `README.md` file further based on your project's specific details and requirements.
```
