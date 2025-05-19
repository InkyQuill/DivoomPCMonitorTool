using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;

namespace DivoomPC.Controls
{
    public class Card : ContentControl
    {
        public new static readonly DependencyProperty ContentProperty =
            DependencyProperty.Register(nameof(Content), typeof(object), typeof(Card), new PropertyMetadata(null));

        public new static readonly DependencyProperty ContentTemplateProperty =
            DependencyProperty.Register(nameof(ContentTemplate), typeof(DataTemplate), typeof(Card), new PropertyMetadata(null));

        public new object Content
        {
            get => GetValue(ContentProperty);
            set => SetValue(ContentProperty, value);
        }

        public new DataTemplate ContentTemplate
        {
            get => (DataTemplate)GetValue(ContentTemplateProperty);
            set => SetValue(ContentTemplateProperty, value);
        }

        public Card()
        {
            DefaultStyleKey = typeof(Card);
        }
    }
} 