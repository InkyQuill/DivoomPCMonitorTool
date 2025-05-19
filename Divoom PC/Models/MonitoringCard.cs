using System.ComponentModel;
using System.Runtime.CompilerServices;

namespace DivoomPC.Models
{
    public class MonitoringCard : INotifyPropertyChanged
    {
        private string _title;
        private string _value;

        public MonitoringCard(string title)
        {
            _title = title;
            _value = "N/A";
        }

        public string Title
        {
            get => _title;
            set
            {
                if (_title != value)
                {
                    _title = value;
                    OnPropertyChanged();
                }
            }
        }

        public string Value
        {
            get => _value;
            set
            {
                if (_value != value)
                {
                    _value = value;
                    OnPropertyChanged();
                }
            }
        }

        public event PropertyChangedEventHandler? PropertyChanged;
        protected virtual void OnPropertyChanged([CallerMemberName] string? propertyName = null)
        {
            PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
        }
    }
} 