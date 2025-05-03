using System;
using System.IO;
using System.Net;
using System.Text;

namespace DivoomPCMonitor
{
    public class HttpClient
    {
        public static int Post(string url, string sendData, out string result)
        {
            result = "";
            try
            {
                byte[] data = Encoding.UTF8.GetBytes(sendData);
                HttpWebRequest wbRequest = (HttpWebRequest)WebRequest.Create(url);
                wbRequest.Proxy = null;
                wbRequest.Method = "POST";
                wbRequest.ContentType = "application/json";
                wbRequest.ContentLength = data.Length;
                wbRequest.Timeout = 1000;

                using (Stream wStream = wbRequest.GetRequestStream())
                {
                    wStream.Write(data, 0, data.Length);
                }

                HttpWebResponse wbResponse = (HttpWebResponse)wbRequest.GetResponse();
                using (Stream responseStream = wbResponse.GetResponseStream())
                {
                    using (StreamReader sReader = new StreamReader(responseStream, Encoding.UTF8))
                    {
                        result = sReader.ReadToEnd();
                    }
                }
            }
            catch (Exception e)
            {
                result = e.Message;
                return -1;
            }
            return 0;
        }

        public static string Get(string url)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(url);
            request.Method = "GET";
            HttpWebResponse response = (HttpWebResponse)request.GetResponse();
            Stream myResponseStream = response.GetResponseStream();
            StreamReader myStreamReader = new StreamReader(myResponseStream, Encoding.GetEncoding("utf-8"));
            string retString = myStreamReader.ReadToEnd();
            myStreamReader.Close();
            myResponseStream.Close();
            return retString;
        }
    }
} 