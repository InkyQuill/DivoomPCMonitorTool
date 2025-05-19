using System;

namespace DivoomPC.Models
{
    public class DivoomDeviceInfo
    {
        public int DeviceId { get; set; }
        public int Hardware { get; set; }
        public required string DeviceName { get; set; }
        public required string DevicePrivateIP { get; set; }
        public required string DeviceMac { get; set; }
    }

    public class DivoomDeviceList
    {
        public required DivoomDeviceInfo[] DeviceList { get; set; }
    }

    public class DivoomTimeGateIndependenceInfo
    {
        public int LcdIndependence { get; set; }
        public int ChannelType { get; set; }
        public int ClockId { get; set; }
    }

    public class DivoomDeviceSelectClockInfo
    {
        public int LcdIndependence { get; set; }
        public int DeviceId { get; set; }
        public int LcdIndex { get; set; }
        public int ClockId { get; set; }
        public required string Command { get; set; }
    }

    public class DivoomDevicePostItem
    {
        public int LcdId { get; set; }
        public required string[] DispData { get; set; }
    }

    public class DivoomDevicePostList
    {
        public required string Command { get; set; }
        public required DivoomDevicePostItem[] ScreenList { get; set; }
    }
} 