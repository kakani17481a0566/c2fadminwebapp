package ServiceApi

import (
	"encoding/json"
	"testing"
)

var req = `{
    "Tksp": {
        "ServiceInterface": {
            "Name": "Login",
            "Version": "1.0",
            "Channel": {
                "Name": "ADC",
                "Type": "WEB",
                "HostIP": {
                    "Type": "IPV4",
                    "Value": "127.0.0.1"
                }
            },
            "Device": {
                "OS": "MAC OS X",
                "Version": "10.9",
                "Id": "65:64:43:54:gg:54",
                "IP": {
                    "Type": "IPV4",
                    "Value": "localhost"
                }
            },
            "TransactionData": {
                "TimeStamp": "2015-09-22 16:45:22",
                "RequestID": "1234567890987654321",
                "TransactionDescription": "B2B Agent login for 919620474150",
                "ServiceDetail": {
                    "BusinessCategory": "B2B",
                    "Partner": "NA",
                    "Category": "LOGIN",
                    "SubCategory": "NA",
                    "Denomination": "0"
                }
            },
            "Users": [
                {
                    "Username": "919620474150",
                    "Type": "AGENT"
                }
             ],
            "Security": {
                "Data": "DRgbfjopwDlouIAuqo0YKWsE+ilW7QyvDsx7M1dbe6IBJB5eexPnhVm9ZBp2/lHE",
                "Skey": "YBpu/nlThSUfd3nBx5Tf1ZkmtJwYsfLodxjpwpn9/pw2HcROcJYxcvM5SxIlNfeHvzVs0CoYLiHRgd7sIuRvW8Z7Pv5fBwKH44/3sMtNNKprC3GEjsy9ur/6JByaTNuxaAKF68JwAvkA6End/NhNbcIcGL7CkCN+R3iwAujBvj9bOi7E91JerLvNRIXAJO10geE6xT12vUfjzxpt5bcr9gjT8Yx1HbBj5mFXi24fR5srECp4qn4WXAyJL2gZ85anj5P2zGaktlaKe7Y/VF1FJekx/7l76KVEoqkkKSx/OCpgodmzMEnlHUtR02n7grK9QV8/9f0/VvGAuh6DboYU1Q==",
                "Hmac": "psqzDKGkYWsQSrrEURtQqMR7Y18GsmY/r+1OXCQ54iABJB5eexPnhVm9ZBp2/lHE"
            }
        }
    }
}`

func TestServiceApi(t *testing.T) {
	var obj_out Root
	err := json.Unmarshal([]byte(req), &obj_out)
	if err != nil {
		t.Error(err)
	}
}
