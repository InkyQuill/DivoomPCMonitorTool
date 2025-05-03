using System;
using System.Windows.Forms;
using System.Threading;
using System.IO;
using Microsoft.Win32;
using Microsoft.Win32.TaskScheduler;
using LibreHardwareMonitor.Hardware;
using System.Globalization;

namespace DivoomPCMonitor
{
    public partial class Form1 : Form
    {
        private const string RegistryKeyPath = @"SOFTWARE\DivoomPCMonitor";
        private const string AppName = "DivoomPCMonitor";
        private const string TaskName = "DivoomPCMonitor";

        private readonly SystemMonitor systemMonitor;
        private readonly DivoomDeviceManager divoomManager;
        private bool isMinimizedToTray = false;

        public Form1()
        {
            // Устанавливаем язык в зависимости от локали системы
            Thread.CurrentThread.CurrentUICulture = CultureInfo.CurrentUICulture;
            
            InitializeComponent();
            this.LCDMsg.Visible = false;
            this.LCDList.Visible = false;
            
            systemMonitor = new SystemMonitor();
            divoomManager = new DivoomDeviceManager();
            
            this.DivoomUpdateDeviceList();
            LoadSettings();

            // Configure tray icon
            notifyIcon.Icon = this.Icon;
            notifyIcon.Text = Properties.Strings.AppTitle;
            notifyIcon.Visible = true;
            notifyIcon.ContextMenuStrip = contextMenuStrip;

            // Add event handlers
            notifyIcon.DoubleClick += notifyIcon_DoubleClick;
            notifyIcon.MouseClick += notifyIcon_MouseClick;
        }

        private void LoadSettings()
        {
            try
            {
                using (RegistryKey key = Registry.CurrentUser.OpenSubKey(RegistryKeyPath))
                {
                    if (key != null)
                    {
                        autoStartCheckBox.Checked = Convert.ToBoolean(key.GetValue("AutoStart", false));
                        startMinimizedCheckBox.Checked = Convert.ToBoolean(key.GetValue("StartMinimized", false));
                        string selectedClock = key.GetValue("SelectedClock", "") as string;
                        if (!string.IsNullOrEmpty(selectedClock))
                        {
                            int index = LCDList.Items.IndexOf(selectedClock);
                            if (index >= 0)
                            {
                                LCDList.SelectedIndex = index;
                            }
                        }
                        isMinimizedToTray = startMinimizedCheckBox.Checked;
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Error loading settings: {ex.Message}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private void SaveSettings()
        {
            try
            {
                using (RegistryKey key = Registry.CurrentUser.CreateSubKey(RegistryKeyPath))
                {
                    key.SetValue("AutoStart", autoStartCheckBox.Checked);
                    key.SetValue("StartMinimized", startMinimizedCheckBox.Checked);
                    if (LCDList.SelectedItem != null)
                    {
                        key.SetValue("SelectedClock", LCDList.SelectedItem.ToString());
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Error saving settings: {ex.Message}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private void SetAutoStart(bool enable)
        {
            try
            {
                using (TaskService ts = new TaskService())
                {
                    if (enable)
                    {
                        // Создаем задачу
                        TaskDefinition td = ts.NewTask();
                        td.RegistrationInfo.Description = "Divoom PC Monitor";
                        td.Principal.RunLevel = TaskRunLevel.Highest;
                        
                        // Настраиваем триггер при входе в систему
                        td.Triggers.Add(new LogonTrigger());
                        
                        // Настраиваем действие
                        td.Actions.Add(new ExecAction(Application.ExecutablePath, null, Path.GetDirectoryName(Application.ExecutablePath)));
                        
                        // Регистрируем задачу
                        ts.RootFolder.RegisterTaskDefinition(TaskName, td);
                    }
                    else
                    {
                        // Удаляем задачу
                        ts.RootFolder.DeleteTask(TaskName, false);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Error configuring autostart: {ex.Message}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private void autoStartCheckBox_CheckedChanged(object sender, EventArgs e)
        {
            SetAutoStart(autoStartCheckBox.Checked);
            SaveSettings();
        }

        private void startMinimizedCheckBox_CheckedChanged(object sender, EventArgs e)
        {
            isMinimizedToTray = startMinimizedCheckBox.Checked;
            SaveSettings();
        }

        private void notifyIcon_DoubleClick(object sender, EventArgs e)
        {
            Show();
            WindowState = FormWindowState.Normal;
            isMinimizedToTray = false;
        }

        private void notifyIcon_MouseClick(object sender, MouseEventArgs e)
        {
            if (e.Button == MouseButtons.Right)
            {
                contextMenuStrip.Show(Cursor.Position);
            }
        }

        private void showMenuItem_Click(object sender, EventArgs e)
        {
            Show();
            WindowState = FormWindowState.Normal;
            isMinimizedToTray = false;
        }

        private void exitMenuItem_Click(object sender, EventArgs e)
        {
            SaveSettings();
            Application.Exit();
        }

        protected override void OnResize(EventArgs e)
        {
            base.OnResize(e);
            if (WindowState == FormWindowState.Minimized)
            {
                Hide();
                isMinimizedToTray = true;
            }
        }

        protected override void OnFormClosing(FormClosingEventArgs e)
        {
            if (e.CloseReason == CloseReason.UserClosing)
            {
                e.Cancel = true;
                WindowState = FormWindowState.Minimized;
                Hide();
                isMinimizedToTray = true;
            }
            else
            {
                timer.Stop();
                timer.Dispose();
                systemMonitor.Close();
            }
            base.OnFormClosing(e);
        }

        protected override void OnLoad(EventArgs e)
        {
            base.OnLoad(e);
            if (isMinimizedToTray)
            {
                WindowState = FormWindowState.Minimized;
                Hide();
            }
        }

        private void DivoomSendHttpInfo(object sender, EventArgs e)
        {
            systemMonitor.Update();

            string cpuTemp = systemMonitor.GetCpuTemperature();
            string cpuUse = systemMonitor.GetCpuUsage();
            string gpuTemp = systemMonitor.GetGpuTemperature();
            string gpuUse = systemMonitor.GetGpuUsage();
            string dispUse = systemMonitor.GetMemoryUsage();
            string hardDiskUse = systemMonitor.GetHardDiskTemperature();

            divoomManager.SendSystemInfo(cpuTemp, cpuUse, gpuTemp, gpuUse, dispUse, hardDiskUse);

            this.CpuTemp.Text = string.Format(Properties.Strings.CpuTemperature, cpuTemp);
            this.CpuUse.Text = string.Format(Properties.Strings.CpuUsage, cpuUse);
            this.GpuUse.Text = string.Format(Properties.Strings.GpuUsage, gpuUse);
            this.GpuTemp.Text = string.Format(Properties.Strings.GpuTemperature, gpuTemp);
            this.HddUse.Text = string.Format(Properties.Strings.HddTemperature, hardDiskUse);
            this.DispUse.Text = string.Format(Properties.Strings.MemoryUsage, dispUse);
        }

        private void DivoomUpdateDeviceList()
        {
            divoomManager.UpdateDeviceList();
            this.divoomList.Items.Clear();
            for (int i = 0; divoomManager.LocalList?.DeviceList != null && i < divoomManager.LocalList.DeviceList.Length; i++)
            {
                this.divoomList.Items.Add(divoomManager.LocalList.DeviceList[i].DeviceName);
            }

            // If only one device is found, select it automatically
            if (divoomManager.LocalList?.DeviceList != null && divoomManager.LocalList.DeviceList.Length == 1)
            {
                divoomList.SelectedIndex = 0;
                
                // After successful device selection, check for saved display in registry
                try
                {
                    using (RegistryKey key = Registry.CurrentUser.OpenSubKey(RegistryKeyPath))
                    {
                        if (key != null)
                        {
                            string selectedClock = key.GetValue("SelectedClock", "") as string;
                            if (!string.IsNullOrEmpty(selectedClock))
                            {
                                int index = LCDList.Items.IndexOf(selectedClock);
                                if (index >= 0)
                                {
                                    LCDList.SelectedIndex = index;
                                }
                            }
                        }
                    }
                }
                catch (Exception ex)
                {
                    MessageBox.Show($"Error loading display settings: {ex.Message}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                }
            }
        }

        private void refreshList_Click(object sender, EventArgs e)
        {
            this.DivoomUpdateDeviceList();
        }

        private void divoomList_SelectedIndexChanged(object sender, EventArgs e)
        {
            if (divoomList.SelectedIndex >= 0)
            {
                divoomManager.SelectDevice(divoomList.SelectedIndex);

                if (divoomManager.LocalList?.DeviceList != null && 
                    divoomList.SelectedIndex < divoomManager.LocalList.DeviceList.Length &&
                    divoomManager.LocalList.DeviceList[divoomList.SelectedIndex].Hardware == 400)
                {
                    this.LCDMsg.Visible = true;
                    this.LCDList.Visible = true;
                    this.LCDList.Enabled = true;
                }
                else
                {
                    this.LCDMsg.Visible = false;
                    this.LCDList.Visible = false;
                    this.LCDList.Enabled = false;
                }
            }
        }

        private void LCDList_SelectedIndexChanged(object sender, EventArgs e)
        {
            if (LCDList.SelectedItem != null)
            {
                divoomManager.SelectLCD(LCDList.SelectedIndex);
            }
        }
    }

    public class DivoomDeviceSelectClockInfo
    {
        public int LcdIndependence { get; set; }
        public int DeviceId { get; set; }
        public int LcdIndex { get; set; }
        public int ClockId { get; set; }
        public string Command { get; set; }
    }
    public class DivoomTimeGateIndependenceInfo
    {
        public int LcdIndependence { get; set; }
        public int ChannelType { get; set; }
        public int ClockId { get; set; }
    }


    public class DivoomDevicePostItem
    {
        public int LcdId { get; set; }


        public string[] DispData { get; set; }

    }
    public class DivoomDevicePostList
    {
        public string Command { get; set; }
        public DivoomDevicePostItem[] ScreenList { get; set; }

    }
    public class UpdateVisitor : IVisitor
    {
        public void VisitComputer(IComputer computer)
        {
            computer.Traverse(this);
        }

        public void VisitHardware(IHardware hardware)
        {
            hardware.Update();
            foreach (IHardware subHardware in hardware.SubHardware)
                subHardware.Accept(this);
        }

        public void VisitSensor(ISensor sensor) { }

        public void VisitParameter(IParameter parameter) { }
    }

}

