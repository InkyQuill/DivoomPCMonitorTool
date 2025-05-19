using Microsoft.UI.Xaml.Controls;
using DivoomPC.ViewModels;

namespace DivoomPC.Views
{
    public sealed partial class SettingsPage : Page
    {
        public MainViewModel ViewModel { get; }

        public SettingsPage()
        {
            this.InitializeComponent();
            ViewModel = new MainViewModel();
        }
    }
} 