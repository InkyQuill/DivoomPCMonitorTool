using Microsoft.UI.Xaml.Controls;
using DivoomPC.ViewModels;

namespace DivoomPC.Views
{
    public sealed partial class MainPage : Page
    {
        public MainViewModel ViewModel { get; }

        public MainPage()
        {
            this.InitializeComponent();
            ViewModel = new MainViewModel();
        }
    }
} 