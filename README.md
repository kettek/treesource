# Init submodules
Acquire the icons:

```
git submodule update --init --recursive
```

# Build GUI

To build the Wails-backed GUI, run:

```
wails generate module
wails build
```

To develop the Wails GUI, run the following commands in separate terminals:

```
wails dev
```

```
cd frontend
npm run dev
```

# Build TUI

To build the TUI, run:

```
go build -tags tui
```
