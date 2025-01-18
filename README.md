# MultiGlass TUI - A Looking Glass Tool

![Screenshot](https://raw.githubusercontent.com/drksbr/lg2/refs/heads/main/screenshot.gif)

## Overview

**MultiGlass TUI** is a terminal-based interface for analyzing BGP (Border Gateway Protocol) peer data. This tool allows network administrators and enthusiasts to view, query, and navigate through peer information in a visually appealing and efficient manner.

The interface is built using the powerful `tview` and `tcell` libraries, providing a responsive and interactive TUI experience.

---

## Features

- **Logo and Shortcuts Display**: A visually distinct top bar showing the application logo and useful shortcuts.
- **Peer List Navigation**: A scrollable list of peers with support for keyboard navigation.
- **Detailed Peer Information**: Display detailed peer data, including AS-PATH and sequence information.
- **Search and Query Modals**: Easily search for peers or initiate new queries using modal dialogs.
- **Keyboard Shortcuts**:
  - `[←]` and `[→]` to navigate between peers.
  - `[Tab]` to cycle focus between components.
  - `[f]` to search for a peer.
  - `[n]` to create a new query.
  - `[q]` to quit the application.

---

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/drksbr/lg2.git
   cd multiglass-tui
   ```

2. **Install Dependencies**:
   Ensure you have Go installed. Then, run:

   ```bash
   go mod tidy
   ```

3. **Build and Run**:
   ```bash
   go run main.go
   ```

---

## Usage

### Starting the Application

Run the application to view the main interface:

```bash
go run main.go
```

### Navigating

- **Select a Peer**: Use `[↓]` and `[↑]` to scroll through the list of peers.
- **Switch Focus**: Press `[Tab]` to toggle focus between the peer list and content pane.
- **Navigate Between Peers**: Use `[←]` and `[→]` to cycle through peer details.
- **Quit**: Press `[q]` to exit the application.

### Search and Query

- **Search for a Peer**: Press `[f]` to open a search modal. Enter the desired peer name and press `[Enter]` to filter the list.
- **Create a New Query**: Press `[n]` to open a query modal. Enter the query details and press `[Enter]`.

---

## Screenshot

![Screenshot](https://via.placeholder.com/800x400?text=MultiGlass+TUI+Screenshot)

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add feature or fix description"
   ```
4. Push your branch:
   ```bash
   git push origin feature-name
   ```
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Acknowledgments

- **[tview](https://github.com/rivo/tview)**: A powerful library for building terminal user interfaces in Go.
- **[tcell](https://github.com/gdamore/tcell)**: A library for handling terminal I/O.

---

## Contact

For questions or support, feel free to reach out:

- Email: isaacramon.adv@gmail.com
- GitHub: [drksbr](https://github.com/drksbr)
