using System;
using Newtonsoft.Json;

namespace DivoomPCMonitor
{
    public class DivoomDeviceInfo
    {
        public int DeviceId { get; set; }

        public int Hardware { get; set; }

        public string DeviceName { get; set; }
        public string DevicePrivateIP { get; set; }
        public string DeviceMac { get; set; }

    }
    public class DivoomDeviceList
    {

        public DivoomDeviceInfo[] DeviceList { get; set; }

    }
    public class DivoomDeviceManager
    {
        private DivoomDeviceList localList;
        private string deviceIPAddr;
        private int selectLCDID = -1;
        private int lcdIndependence;

        public DivoomDeviceList LocalList => localList;
        public string DeviceIPAddr => deviceIPAddr;
        public int SelectLCDID => selectLCDID;
        public int LcdIndependence => lcdIndependence;

        public DivoomDeviceManager()
        {
        
        }

        public void UpdateDeviceList()
        {
            string url_info = "http://app.divoom-gz.com/Device/ReturnSameLANDevice";
            string device_list = HttpClient.Get(url_info);
            localList = JsonConvert.DeserializeObject<DivoomDeviceList>(device_list);
        }

        public void SelectDevice(int index)
        {
            if (index >= 0 && index < localList.DeviceList.Length)
            {
                deviceIPAddr = localList.DeviceList[index].DevicePrivateIP;

                if (localList.DeviceList[index].Hardware == 400)
                {
                    string url_info = "http://app.divoom-gz.com/Channel/Get5LcdInfoV2?DeviceType=LCD&DeviceId=" + 
                                    localList.DeviceList[index].DeviceId;
                    string IndependenceStr = HttpClient.Get(url_info);
                    if (!string.IsNullOrEmpty(IndependenceStr))
                    {
                        DivoomTimeGateIndependenceInfo IndependenceInfo = 
                            JsonConvert.DeserializeObject<DivoomTimeGateIndependenceInfo>(IndependenceStr);
                        lcdIndependence = IndependenceInfo.LcdIndependence;
                    }
                }
                else
                {
                    selectLCDID = 0;
                }
            }
        }

        public void SelectLCD(int index)
        {
            if (index >= 0)
            {
                selectLCDID = index - 1;

                if (localList != null && localList.DeviceList != null && localList.DeviceList.Length > 0)
                {
                    DivoomDeviceSelectClockInfo PostInfo = new DivoomDeviceSelectClockInfo
                    {
                        LcdIndependence = lcdIndependence,
                        Command = "Channel/SetClockSelectId",
                        LcdIndex = index,
                        ClockId = 625
                    };

                    string para_info = JsonConvert.SerializeObject(PostInfo);
                    string response_info;
                    HttpClient.Post("http://" + deviceIPAddr + ":80/post", para_info, out response_info);
                }
            }
        }

        public void SendSystemInfo(string cpuTemp, string cpuUse, string gpuTemp, string gpuUse, 
                                 string dispUse, string hardDiskUse)
        {
            if (string.IsNullOrEmpty(deviceIPAddr) || localList?.DeviceList == null || 
                localList.DeviceList.Length == 0 || selectLCDID < 0)
            {
                return;
            }

            DivoomDevicePostList PostInfo = new DivoomDevicePostList();
            DivoomDevicePostItem PostItem = new DivoomDevicePostItem
            {
                LcdId = selectLCDID,
                DispData = new string[6]
            };

            PostItem.DispData[2] = cpuTemp;
            PostItem.DispData[0] = cpuUse;
            PostItem.DispData[3] = gpuTemp;
            PostItem.DispData[1] = gpuUse;
            PostItem.DispData[5] = hardDiskUse;
            PostItem.DispData[4] = dispUse;

            PostInfo.Command = "Device/UpdatePCParaInfo";
            PostInfo.ScreenList = new DivoomDevicePostItem[] { PostItem };

            string para_info = JsonConvert.SerializeObject(PostInfo);
            string response_info;
            HttpClient.Post("http://" + deviceIPAddr + ":80/post", para_info, out response_info);
        }
    }
} 