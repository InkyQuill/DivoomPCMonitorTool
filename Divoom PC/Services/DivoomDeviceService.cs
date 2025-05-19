using System;
using System.Threading.Tasks;
using DivoomPC.Models;
using System.Net.Http;
using Newtonsoft.Json;
using Windows.Storage;

namespace DivoomPC.Services
{
    public class DivoomDeviceService
    {
        private readonly HttpClient _httpClient;
        private DivoomDeviceList? _localList;
        private string? _deviceIPAddr;
        private int _selectLCDID = -1;
        private int _lcdIndependence;

        public DivoomDeviceList? LocalList => _localList;
        public string? DeviceIPAddr => _deviceIPAddr;
        public int SelectLCDID => _selectLCDID;
        public int LcdIndependence => _lcdIndependence;

        public DivoomDeviceService()
        {
            _httpClient = new HttpClient();
        }

        public async Task UpdateDeviceListAsync()
        {
            string url_info = "http://app.divoom-gz.com/Device/ReturnSameLANDevice";
            string device_list = await _httpClient.GetStringAsync(url_info);
            _localList = JsonConvert.DeserializeObject<DivoomDeviceList>(device_list);
        }

        public async Task SelectDeviceAsync(int index)
        {
            if (index >= 0 && _localList?.DeviceList != null && index < _localList.DeviceList.Length)
            {
                _deviceIPAddr = _localList.DeviceList[index].DevicePrivateIP;

                if (_localList.DeviceList[index].Hardware == 400)
                {
                    string url_info = "http://app.divoom-gz.com/Channel/Get5LcdInfoV2?DeviceType=LCD&DeviceId=" + 
                                    _localList.DeviceList[index].DeviceId;
                    string IndependenceStr = await _httpClient.GetStringAsync(url_info);
                    if (!string.IsNullOrEmpty(IndependenceStr))
                    {
                        DivoomTimeGateIndependenceInfo? IndependenceInfo = 
                            JsonConvert.DeserializeObject<DivoomTimeGateIndependenceInfo>(IndependenceStr);
                        if (IndependenceInfo != null)
                        {
                            _lcdIndependence = IndependenceInfo.LcdIndependence;
                        }
                    }
                }
                else
                {
                    _selectLCDID = 0;
                }
            }
        }

        public async Task SelectLCDAsync(int index)
        {
            if (index >= 0)
            {
                _selectLCDID = index - 1;

                if (_localList?.DeviceList != null && _localList.DeviceList.Length > 0)
                {
                    DivoomDeviceSelectClockInfo PostInfo = new()
                    {
                        LcdIndependence = _lcdIndependence,
                        Command = "Channel/SetClockSelectId",
                        LcdIndex = index,
                        ClockId = 625
                    };

                    string para_info = JsonConvert.SerializeObject(PostInfo);
                    var content = new StringContent(para_info);
                    await _httpClient.PostAsync($"http://{_deviceIPAddr}:80/post", content);
                }
            }
        }

        public async Task SendSystemInfoAsync(string cpuTemp, string cpuUse, string gpuTemp, string gpuUse, 
                                           string dispUse, string hardDiskUse)
        {
            if (string.IsNullOrEmpty(_deviceIPAddr) || _localList?.DeviceList == null || 
                _localList.DeviceList.Length == 0 || _selectLCDID < 0)
            {
                return;
            }

            DivoomDevicePostItem PostItem = new()
            {
                LcdId = _selectLCDID,
                DispData = new string[6]
            };

            PostItem.DispData[2] = cpuTemp;
            PostItem.DispData[0] = cpuUse;
            PostItem.DispData[3] = gpuTemp;
            PostItem.DispData[1] = gpuUse;
            PostItem.DispData[5] = hardDiskUse;
            PostItem.DispData[4] = dispUse;

            DivoomDevicePostList PostInfo = new()
            {
                Command = "Device/UpdatePCParaInfo",
                ScreenList = new[] { PostItem }
            };

            string para_info = JsonConvert.SerializeObject(PostInfo);
            var content = new StringContent(para_info);
            await _httpClient.PostAsync($"http://{_deviceIPAddr}:80/post", content);
        }
    }
} 