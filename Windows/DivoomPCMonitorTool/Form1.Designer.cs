namespace DivoomPCMonitor
{
    partial class Form1
    {
        /// <summary>
        /// 必需的设计器变量。
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// 清理所有正在使用的资源。
        /// </summary>
        /// <param name="disposing">如果应释放托管资源，为 true；否则为 false。</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows 窗体设计器生成的代码

        /// <summary>
        /// 设计器支持所需的方法 - 不要
        /// 使用代码编辑器修改此方法的内容。
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Form1));
            this.refreshList = new System.Windows.Forms.Button();
            this.CpuUse = new System.Windows.Forms.Label();
            this.CpuTemp = new System.Windows.Forms.Label();
            this.GpuUse = new System.Windows.Forms.Label();
            this.GpuTemp = new System.Windows.Forms.Label();
            this.DispUse = new System.Windows.Forms.Label();
            this.HddUse = new System.Windows.Forms.Label();
            this.divoomList = new System.Windows.Forms.ListBox();
            this.LCDList = new System.Windows.Forms.ListBox();
            this.LCDMsg = new System.Windows.Forms.Label();
            this.DeviceListMsg = new System.Windows.Forms.Label();
            this.HardwareInfo = new System.Windows.Forms.Label();
            this.notifyIcon = new System.Windows.Forms.NotifyIcon(this.components);
            this.contextMenuStrip = new System.Windows.Forms.ContextMenuStrip(this.components);
            this.showMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.exitMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.autoStartCheckBox = new System.Windows.Forms.CheckBox();
            this.startMinimizedCheckBox = new System.Windows.Forms.CheckBox();
            this.timer = new System.Windows.Forms.Timer(this.components);
            this.contextMenuStrip.SuspendLayout();
            this.SuspendLayout();
            // 
            // refreshList
            // 
            this.refreshList.Location = new System.Drawing.Point(16, 316);
            this.refreshList.Margin = new System.Windows.Forms.Padding(4);
            this.refreshList.Name = "refreshList";
            this.refreshList.Size = new System.Drawing.Size(173, 31);
            this.refreshList.TabIndex = 0;
            this.refreshList.Text = global::DivoomPCMonitor.Properties.Strings.Form1_RefreshList;
            this.refreshList.UseVisualStyleBackColor = true;
            this.refreshList.Click += new System.EventHandler(this.refreshList_Click);
            // 
            // CpuUse
            // 
            this.CpuUse.AutoSize = true;
            this.CpuUse.Location = new System.Drawing.Point(350, 80);
            this.CpuUse.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.CpuUse.Name = "CpuUse";
            this.CpuUse.Size = new System.Drawing.Size(0, 16);
            this.CpuUse.TabIndex = 1;
            // 
            // CpuTemp
            // 
            this.CpuTemp.AutoSize = true;
            this.CpuTemp.Location = new System.Drawing.Point(350, 120);
            this.CpuTemp.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.CpuTemp.Name = "CpuTemp";
            this.CpuTemp.Size = new System.Drawing.Size(0, 16);
            this.CpuTemp.TabIndex = 1;
            // 
            // GpuUse
            // 
            this.GpuUse.AutoSize = true;
            this.GpuUse.Location = new System.Drawing.Point(350, 160);
            this.GpuUse.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.GpuUse.Name = "GpuUse";
            this.GpuUse.Size = new System.Drawing.Size(0, 16);
            this.GpuUse.TabIndex = 1;
            // 
            // GpuTemp
            // 
            this.GpuTemp.AutoSize = true;
            this.GpuTemp.Location = new System.Drawing.Point(350, 200);
            this.GpuTemp.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.GpuTemp.Name = "GpuTemp";
            this.GpuTemp.Size = new System.Drawing.Size(0, 16);
            this.GpuTemp.TabIndex = 1;
            // 
            // DispUse
            // 
            this.DispUse.AutoSize = true;
            this.DispUse.Location = new System.Drawing.Point(350, 240);
            this.DispUse.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.DispUse.Name = "DispUse";
            this.DispUse.Size = new System.Drawing.Size(0, 16);
            this.DispUse.TabIndex = 1;
            // 
            // HddUse
            // 
            this.HddUse.AutoSize = true;
            this.HddUse.Location = new System.Drawing.Point(350, 280);
            this.HddUse.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.HddUse.Name = "HddUse";
            this.HddUse.Size = new System.Drawing.Size(0, 16);
            this.HddUse.TabIndex = 1;
            // 
            // divoomList
            // 
            this.divoomList.ItemHeight = 16;
            this.divoomList.Location = new System.Drawing.Point(16, 80);
            this.divoomList.Margin = new System.Windows.Forms.Padding(4);
            this.divoomList.Name = "divoomList";
            this.divoomList.Size = new System.Drawing.Size(173, 228);
            this.divoomList.TabIndex = 2;
            this.divoomList.SelectedIndexChanged += new System.EventHandler(this.divoomList_SelectedIndexChanged);
            // 
            // LCDList
            // 
            this.LCDList.Enabled = false;
            this.LCDList.ItemHeight = 16;
            this.LCDList.Items.AddRange(new object[] {
            global::DivoomPCMonitor.Properties.Strings.NotSelected,
            "1",
            "2",
            "3",
            "4",
            "5"});
            this.LCDList.Location = new System.Drawing.Point(197, 80);
            this.LCDList.Margin = new System.Windows.Forms.Padding(4);
            this.LCDList.Name = "LCDList";
            this.LCDList.Size = new System.Drawing.Size(145, 276);
            this.LCDList.TabIndex = 2;
            this.LCDList.SelectedIndexChanged += new System.EventHandler(this.LCDList_SelectedIndexChanged);
            // 
            // LCDMsg
            // 
            this.LCDMsg.Location = new System.Drawing.Point(197, 53);
            this.LCDMsg.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.LCDMsg.Name = "LCDMsg";
            this.LCDMsg.Size = new System.Drawing.Size(145, 27);
            this.LCDMsg.TabIndex = 3;
            this.LCDMsg.Text = Properties.Strings.Form1_LCDMsg;
            // 
            // DeviceListMsg
            // 
            this.DeviceListMsg.Location = new System.Drawing.Point(16, 53);
            this.DeviceListMsg.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.DeviceListMsg.Name = "DeviceListMsg";
            this.DeviceListMsg.Size = new System.Drawing.Size(133, 27);
            this.DeviceListMsg.TabIndex = 4;
            this.DeviceListMsg.Text = Properties.Strings.Form1_DeviceListMsg;
            // 
            // HardwareInfo
            // 
            this.HardwareInfo.Location = new System.Drawing.Point(350, 53);
            this.HardwareInfo.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.HardwareInfo.Name = "HardwareInfo";
            this.HardwareInfo.Size = new System.Drawing.Size(267, 27);
            this.HardwareInfo.TabIndex = 5;
            this.HardwareInfo.Text = Properties.Strings.Form1_HardwareInfo;
            // 
            // notifyIcon
            // 
            this.notifyIcon.ContextMenuStrip = this.contextMenuStrip;
            this.notifyIcon.Icon = ((System.Drawing.Icon)(resources.GetObject("notifyIcon.Icon")));
            this.notifyIcon.Text = global::DivoomPCMonitor.Properties.Strings.AppTitle;
            this.notifyIcon.Visible = true;
            this.notifyIcon.DoubleClick += new System.EventHandler(this.notifyIcon_DoubleClick);
            // 
            // contextMenuStrip
            // 
            this.contextMenuStrip.ImageScalingSize = new System.Drawing.Size(20, 20);
            this.contextMenuStrip.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.showMenuItem,
            this.exitMenuItem});
            this.contextMenuStrip.Name = "contextMenuStrip";
            this.contextMenuStrip.Size = new System.Drawing.Size(115, 52);
            // 
            // showMenuItem
            // 
            this.showMenuItem.Name = "showMenuItem";
            this.showMenuItem.Size = new System.Drawing.Size(114, 24);
            this.showMenuItem.Text = global::DivoomPCMonitor.Properties.Strings.Form1_ShowMenuItem;
            this.showMenuItem.Click += new System.EventHandler(this.showMenuItem_Click);
            // 
            // exitMenuItem
            // 
            this.exitMenuItem.Name = "exitMenuItem";
            this.exitMenuItem.Size = new System.Drawing.Size(114, 24);
            this.exitMenuItem.Text = global::DivoomPCMonitor.Properties.Strings.Form1_ExitMenuItem;
            this.exitMenuItem.Click += new System.EventHandler(this.exitMenuItem_Click);
            // 
            // autoStartCheckBox
            // 
            this.autoStartCheckBox.AutoSize = true;
            this.autoStartCheckBox.Location = new System.Drawing.Point(16, 12);
            this.autoStartCheckBox.Name = "autoStartCheckBox";
            this.autoStartCheckBox.Size = new System.Drawing.Size(139, 20);
            this.autoStartCheckBox.TabIndex = 6;
            this.autoStartCheckBox.Text = global::DivoomPCMonitor.Properties.Strings.Form1_AutoStartCheckBox;
            this.autoStartCheckBox.UseVisualStyleBackColor = true;
            this.autoStartCheckBox.CheckedChanged += new System.EventHandler(this.autoStartCheckBox_CheckedChanged);
            // 
            // startMinimizedCheckBox
            // 
            this.startMinimizedCheckBox.AutoSize = true;
            this.startMinimizedCheckBox.Location = new System.Drawing.Point(271, 12);
            this.startMinimizedCheckBox.Name = "startMinimizedCheckBox";
            this.startMinimizedCheckBox.Size = new System.Drawing.Size(119, 20);
            this.startMinimizedCheckBox.TabIndex = 7;
            this.startMinimizedCheckBox.Text = global::DivoomPCMonitor.Properties.Strings.Form1_StartMinimizedCheckBox;
            this.startMinimizedCheckBox.UseVisualStyleBackColor = true;
            this.startMinimizedCheckBox.CheckedChanged += new System.EventHandler(this.startMinimizedCheckBox_CheckedChanged);
            // 
            // timer
            // 
            this.timer.Enabled = true;
            this.timer.Interval = 1000;
            this.timer.Tick += new System.EventHandler(this.DivoomSendHttpInfo);
            // 
            // Form1
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(8F, 16F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(619, 371);
            this.Controls.Add(this.CpuUse);
            this.Controls.Add(this.CpuTemp);
            this.Controls.Add(this.GpuUse);
            this.Controls.Add(this.GpuTemp);
            this.Controls.Add(this.DispUse);
            this.Controls.Add(this.HddUse);
            this.Controls.Add(this.refreshList);
            this.Controls.Add(this.divoomList);
            this.Controls.Add(this.LCDList);
            this.Controls.Add(this.LCDMsg);
            this.Controls.Add(this.DeviceListMsg);
            this.Controls.Add(this.HardwareInfo);
            this.Controls.Add(this.autoStartCheckBox);
            this.Controls.Add(this.startMinimizedCheckBox);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedSingle;
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.Margin = new System.Windows.Forms.Padding(4);
            this.MaximizeBox = false;
            this.Name = "Form1";
            this.Text = Properties.Strings.AppTitle;
            this.contextMenuStrip.ResumeLayout(false);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.ListBox divoomList;
        private System.Windows.Forms.ListBox LCDList;
        private System.Windows.Forms.Button refreshList;
        private System.Windows.Forms.Label CpuUse;
        private System.Windows.Forms.Label CpuTemp;
        private System.Windows.Forms.Label GpuUse;
        private System.Windows.Forms.Label GpuTemp;
        private System.Windows.Forms.Label DispUse;
        private System.Windows.Forms.Label HddUse;
        private System.Windows.Forms.Label LCDMsg;
        private System.Windows.Forms.Label DeviceListMsg;
        private System.Windows.Forms.Label HardwareInfo;
        private System.Windows.Forms.Timer timer;

        private System.Windows.Forms.NotifyIcon notifyIcon;
        private System.Windows.Forms.ContextMenuStrip contextMenuStrip;
        private System.Windows.Forms.ToolStripMenuItem showMenuItem;
        private System.Windows.Forms.ToolStripMenuItem exitMenuItem;
        private System.Windows.Forms.CheckBox autoStartCheckBox;
        private System.Windows.Forms.CheckBox startMinimizedCheckBox;

    }
}

