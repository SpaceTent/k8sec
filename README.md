# k8sec

`k8sec`, Simple TUI-based Kubernetes Secret editor. It provides a user-friendly interface for listing, filtering, editing, and adding secrets in your Kubernetes cluster, all from the terminal.

## Features

-   **Interactive List**: Browse all secrets in a namespace with a searchable TUI.
-   **Fuzzy Filtering**: Quickly find the secret you need by typing its name.
-   **Safe Editing**: Opens your preferred editor (e.g., `vim`, `nano`) to edit secrets in a `.env` format.
-   **Change Detection**: Automatically detects if any changes were made. If no changes are detected, it skips the cluster update to avoid unnecessary operations.
-   **Add New Secrets**: Easily create new secrets.

## Prerequisites

-   A running Kubernetes cluster.
-   A valid `~/.kube/config` file.
-   The `EDITOR` environment variable set to your preferred text editor (defaults to `nano`).

## Installation

To build the project, you need Go installed.

```bash
# Clone the repository (if applicable)
# git clone <repository-url>
# cd k8sec

# Build the binary using the Makefile
make build
```

## Usage

Run the binary to start the interactive session:

```bash
# Edit secrets in the default namespace
./k8sec

# Edit secrets in a specific namespace using flags (recommended)
./k8sec -n my-namespace

# Edit secrets in a specific namespace using environment variables
NAMESPACE=my-namespace ./k8sec
```

### Keybindings

-   **Up/Down / j/k**: Navigate the list of secrets.
-   **Enter**: Select the highlighted secret for editing.
-   **/**: Enter filter mode to search for a specific secret.
-   **a**: Add a new secret to the cluster.
-   **Esc**: Clear the filter or go back.
-   **Ctrl+C**: Quit the application.

## How it Works

1.  **Selection**: The tool lists all secrets in the target namespace using `client-go`.
2.  **TUI**: You select a secret using the Bubble Tea-powered list component.
3.  **Editing**: 
    -   The secret's data is decrypted and converted into a temporary `.env` style file (`KEY=VALUE`).
    -   Your system editor is opened with this temporary file.
4.  **Verification**: After saving and closing the editor, `k8sec` compares the original data with the updated file.
5.  **Update**: If changes were made, the secret is updated in the Kubernetes cluster. If no changes were detected, it prints "no changes" and exits.
6.  **Creation**: When adding a new secret, you are prompted for a name and namespace, then dropped into the editor to define the secret's data.

## License

This project is licensed under the MIT Licence – see the [LICENCE](LICENCE) file for details.
