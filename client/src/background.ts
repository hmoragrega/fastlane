"use strict";

import { app, protocol, BrowserWindow, Tray, ipcMain } from "electron";
import { createProtocol } from "vue-cli-plugin-electron-builder/lib";
import installExtension, { VUEJS_DEVTOOLS } from "electron-devtools-installer";

const notifier = require('node-notifier');
notifier.options.customPath = "terminal-notifier"

ipcMain.on('asynchronous-message', (event, arg) => {
  console.log("async", arg) // prints "ping"
  event.reply('asynchronous-reply', 'pong')
})

ipcMain.on('synchronous-message', (event, arg) => {
  console.log("sync", arg) // prints "ping"

  console.log("notification icon", path.join(__dirname, 'img/icons/git-24x24.png'))
  console.log("notification icon", path.join(__static, 'img/icons/git-24x24.png'))

  notifier.notify({
    title: arg.title,
    message: arg.message,
    icon: path.join(__static, 'img/icons/git-24x24.png'),
    sound: true,
    wait: true,
  })
})

const path = require('path')
const url = require('url')

const isDevelopment = process.env.NODE_ENV !== "production";
const assetsDirectory = path.join(__dirname, 'assets')

// Scheme must be registered before the app is ready
protocol.registerSchemesAsPrivileged([
  { scheme: "app", privileges: { secure: true, standard: true } },
]);

let tray: Tray
let win: BrowserWindow

declare const __static: string;

async function createWindow() {
  // Create the browser window.
  win = new BrowserWindow({
    width: 800,
    height: 800,
    show: false,
    frame: false,
    fullscreenable: false,
    resizable: false,
    webPreferences: {
      // Required for Spectron testing
      enableRemoteModule: !!process.env.IS_TEST,

      // eslint-disable-next-line
      preload: path.resolve(__static, 'preload.js'),

      // Use pluginOptions.nodeIntegration, leave this alone
      // See nklayman.github.io/vue-cli-plugin-electron-builder/guide/security.html#node-integration for more info
      nodeIntegration: process.env
        .ELECTRON_NODE_INTEGRATION as unknown as boolean,
      contextIsolation: !process.env.ELECTRON_NODE_INTEGRATION,
      /*nodeIntegration: true,
      contextIsolation: false,
      enableRemoteModule: true,*/
    },
  });

  win.setMaximumSize(800, 1000)

  win.webContents.on('new-window', function(e, url) {
    e.preventDefault();
    require('electron').shell.openExternal(url);
  });

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    console.log("LOADING WEB 1 = ", process.env.WEBPACK_DEV_SERVER_URL)

    // Load the url of the dev server if in development mode
    await win.loadURL(process.env.WEBPACK_DEV_SERVER_URL as string);
    if (!process.env.IS_TEST) win.webContents.openDevTools();
  } else {
    createProtocol("app")
    win.loadURL("app://./index.html");
  }
}

const openDevTools = () => {
  win.webContents.openDevTools();
}

async function createTray() {
  console.log("ASSETS!", assetsDirectory)
  console.log("STATIC!", __static)

  tray = new Tray(path.join(__static, 'img/icons/git-16x16.png'))
  tray.on('right-click', openDevTools)
  tray.on('double-click', toggleWindow)
  tray.on('click', function (event) {
    toggleWindow()

    /*// Show devtools when command clicked
    if (win.isVisible() && process.defaultApp && event.metaKey) {
      win.openDevTools({mode: 'detach'})
    }*/
  })
}


const toggleWindow = () => {
  if (win.isVisible()) {
    win.hide()
  } else {
    showWindow()
  }
}

const showWindow = () => {
  const position = getWindowPosition()
  win.setPosition(position.x, position.y, false)
  win.show()
  win.focus()
}

const getWindowPosition = () => {
  const windowBounds = win.getBounds()
  const trayBounds = tray.getBounds()

  // Center window horizontally below the tray icon
  const x = Math.round(trayBounds.x + (trayBounds.width / 2) - (windowBounds.width / 2))

  // Position window 4 pixels vertically below the tray icon
  const y = Math.round(trayBounds.y + trayBounds.height + 4)

  return {x: x, y: y}
}

// Quit when all windows are closed.
app.on("window-all-closed", () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("activate", () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) createWindow();
});

app.setAppUserModelId("com.hmoragrega.fastlane");

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", async () => {
  if (isDevelopment && !process.env.IS_TEST) {
    // Install Vue Devtools
    try {
      await installExtension(VUEJS_DEVTOOLS);
    } catch (e) {
      console.error("Vue Devtools failed to install:", e.toString());
    }
  }
  createTray();
  createWindow();
  //showWindow();

// Object
/*
  console.log("calling notify")
  notifier.notify({
    title: "Fastlane notification",
    message: "Merge available!",
    customPath: "terminal-notifier"
  })
  console.log("notify called")
*/
});

// Exit cleanly on request from parent process in development mode.
if (isDevelopment) {
  if (process.platform === "win32") {
    process.on("message", (data) => {
      if (data === "graceful-exit") {
        app.quit();
      }
    });
  } else {
    process.on("SIGTERM", () => {
      app.quit();
    });
  }
}
