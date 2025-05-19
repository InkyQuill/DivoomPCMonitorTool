using Microsoft.UI.Xaml;
using Microsoft.UI.Windowing;
using Microsoft.UI;
using System;

namespace DivoomPC.Controls
{
    public class WindowEx : Window
    {
        private AppWindow? _appWindow;

        public WindowEx()
        {
            _appWindow = GetAppWindowForCurrentWindow(this);
        }

        private AppWindow GetAppWindowForCurrentWindow(Window window)
        {
            var windowHandle = WinRT.Interop.WindowNative.GetWindowHandle(window);
            var windowId = Win32Interop.GetWindowIdFromWindow(windowHandle);
            return AppWindow.GetFromWindowId(windowId);
        }

        public void SetTitle(string title)
        {
            if (_appWindow != null)
            {
                _appWindow.Title = title;
            }
        }

        public void SetIcon(string iconPath)
        {
            if (_appWindow != null)
            {
                _appWindow.SetIcon(iconPath);
            }
        }
    }
} 