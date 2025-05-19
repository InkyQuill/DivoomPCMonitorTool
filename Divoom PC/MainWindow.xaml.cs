using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices.WindowsRuntime;
using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Controls.Primitives;
using Microsoft.UI.Xaml.Data;
using Microsoft.UI.Xaml.Input;
using Microsoft.UI.Xaml.Media;
using Microsoft.UI.Xaml.Navigation;
using Windows.Foundation;
using Windows.Foundation.Collections;
using DivoomPC.Views;
using Microsoft.UI;
using Microsoft.UI.Windowing;
using WinRT.Interop;
using System.Runtime.InteropServices;
using H.NotifyIcon;
using System.Windows.Input;
using Microsoft.UI.Xaml.Media.Imaging;
using DivoomPC.Controls;
using DivoomPC.ViewModels;

// To learn more about WinUI, the WinUI project structure,
// and more about our project templates, see: http://aka.ms/winui-project-info.

namespace DivoomPC
{
    /// <summary>
    /// An empty window that can be used on its own or navigated to within a Frame.
    /// </summary>
    public sealed partial class MainWindow : WindowEx
    {
        private readonly MainViewModel _viewModel;
        public ICommand RestoreCommand { get; }
        public ICommand ExitCommand { get; }

        public MainWindow()
        {
            this.InitializeComponent();
            _viewModel = new MainViewModel();
            NavView.DataContext = _viewModel;

            RestoreCommand = new RelayCommand(RestoreFromTray);
            ExitCommand = new RelayCommand(Exit);
            
            NavView.SelectedItem = NavView.MenuItems[0];
            ContentFrame.Navigate(typeof(MainPage));
        }

        private void NavView_SelectionChanged(NavigationView sender, NavigationViewSelectionChangedEventArgs args)
        {
            if (args.SelectedItemContainer is NavigationViewItem item)
            {
                switch (item.Tag.ToString())
                {
                    case "home":
                        ContentFrame.Navigate(typeof(MainPage));
                        break;
                    case "settings":
                        ContentFrame.Navigate(typeof(SettingsPage));
                        break;
                }
            }
        }

        private void MinimizeButton_Click(object sender, RoutedEventArgs e)
        {
            MinimizeToTray();
        }

        private void MinimizeToTray()
        {
            _viewModel.IsMinimizedToTray = true;
            this.Hide();
        }

        private void RestoreFromTray()
        {
            _viewModel.IsMinimizedToTray = false;
            this.Show();
            this.Activate();
        }

        private void Exit()
        {
            Application.Current.Exit();
        }
    }

    public class RelayCommand : ICommand
    {
        private readonly Action _execute;
        private readonly Func<bool>? _canExecute;

        public RelayCommand(Action execute, Func<bool>? canExecute = null)
        {
            _execute = execute ?? throw new ArgumentNullException(nameof(execute));
            _canExecute = canExecute;
        }

        public event EventHandler? CanExecuteChanged;

        public bool CanExecute(object? parameter) => _canExecute?.Invoke() ?? true;

        public void Execute(object? parameter) => _execute();
    }
}
