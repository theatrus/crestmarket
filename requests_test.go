package crestmarket

import (
	"testing"
)

var rootData = []byte(`{"motd": {"dust": {"href": "http://newsfeed.eveonline.com/articles/71"}, "eve": {"href": "http://client.eveonline.com/motd/"}, "server": {"href": "http://client.eveonline.com/motd/"}}, "crestEndpoint": {"href": "https://api-sisi.testeveonline.com/"}, "corporationRoles": {"href": "https://api-sisi.testeveonline.com/corporations/roles/"}, "itemGroups": {"href": "https://api-sisi.testeveonline.com/inventory/groups/"}, "channels": {"href": "https://api-sisi.testeveonline.com/chat/channels/"}, "corporations": {"href": "https://api-sisi.testeveonline.com/corporations/"}, "alliances": {"href": "https://api-sisi.testeveonline.com/alliances/"}, "itemTypes": {"href": "https://api-sisi.testeveonline.com/inventory/types/"}, "decode": {"href": "https://api-sisi.testeveonline.com/decode/"}, "battleTheatres": {"href": "https://api-sisi.testeveonline.com/battles/theatres/"}, "marketPrices": {"href": "https://api-sisi.testeveonline.com/market/prices/"}, "itemCategories": {"href": "https://api-sisi.testeveonline.com/inventory/categories/"}, "regions": {"href": "https://api-sisi.testeveonline.com/regions/"}, "bloodlines": {"href": "https://api-sisi.testeveonline.com/bloodlines/"}, "marketGroups": {"href": "https://api-sisi.testeveonline.com/market/groups/"}, "tournaments": {"href": "https://api-sisi.testeveonline.com/tournaments/"}, "map": {"href": "https://api-sisi.testeveonline.com/map/"}, "virtualGoodStore": {"href": "https://sisivgs.testeveonline.com/"}, "serverVersion": "EVE-2014-ISAIA 8.5.848581.848577", "wars": {"href": "https://api-sisi.testeveonline.com/wars/"}, "incursions": {"href": "https://api-sisi.testeveonline.com/incursions/"}, "races": {"href": "https://api-sisi.testeveonline.com/races/"}, "authEndpoint": {"href": "https://sisilogin.testeveonline.com/oauth/token/"}, "serviceStatus": {"dust": "vip", "eve": "online", "server": "online"}, "userCounts": {"dust": 0, "dust_str": "0", "eve": 114, "eve_str": "114"}, "industry": {"facilities": {"href": "https://api-sisi.testeveonline.com/industry/facilities/"}, "specialities": {"href": "https://api-sisi.testeveonline.com/industry/specialities/"}, "teamsInAuction": {"href": "https://api-sisi.testeveonline.com/industry/teams/auction/"}, "systems": {"href": "https://api-sisi.testeveonline.com/industry/systems/"}, "teams": {"href": "https://api-sisi.testeveonline.com/industry/teams/"}}, "clients": {"dust": {"href": "https://api-sisi.testeveonline.com/roots/dust/"}, "eve": {"href": "https://api-sisi.testeveonline.com/roots/eve/"}}, "time": {"href": "https://api-sisi.testeveonline.com/time/"}, "marketTypes": {"href": "https://api-sisi.testeveonline.com/market/types/"}, "serverName": "SINGULARITY"}`)

var regionsData = []byte(`{"totalCount_str": "99", "items": [{"href": "https://api-sisi.testeveonline.com/regions/11000001/", "name": "A-R00001"}, {"href": "https://api-sisi.testeveonline.com/regions/11000002/", "name": "A-R00002"}, {"href": "https://api-sisi.testeveonline.com/regions/11000003/", "name": "A-R00003"}, {"href": "https://api-sisi.testeveonline.com/regions/10000019/", "name": "A821-A"}, {"href": "https://api-sisi.testeveonline.com/regions/10000054/", "name": "Aridia"}, {"href": "https://api-sisi.testeveonline.com/regions/11000004/", "name": "B-R00004"}, {"href": "https://api-sisi.testeveonline.com/regions/11000005/", "name": "B-R00005"}, {"href": "https://api-sisi.testeveonline.com/regions/11000006/", "name": "B-R00006"}, {"href": "https://api-sisi.testeveonline.com/regions/11000007/", "name": "B-R00007"}, {"href": "https://api-sisi.testeveonline.com/regions/11000008/", "name": "B-R00008"}, {"href": "https://api-sisi.testeveonline.com/regions/10000069/", "name": "Black Rise"}, {"href": "https://api-sisi.testeveonline.com/regions/10000055/", "name": "Branch"}, {"href": "https://api-sisi.testeveonline.com/regions/11000009/", "name": "C-R00009"}, {"href": "https://api-sisi.testeveonline.com/regions/11000010/", "name": "C-R00010"}, {"href": "https://api-sisi.testeveonline.com/regions/11000011/", "name": "C-R00011"}, {"href": "https://api-sisi.testeveonline.com/regions/11000012/", "name": "C-R00012"}, {"href": "https://api-sisi.testeveonline.com/regions/11000013/", "name": "C-R00013"}, {"href": "https://api-sisi.testeveonline.com/regions/11000014/", "name": "C-R00014"}, {"href": "https://api-sisi.testeveonline.com/regions/11000015/", "name": "C-R00015"}, {"href": "https://api-sisi.testeveonline.com/regions/10000007/", "name": "Cache"}, {"href": "https://api-sisi.testeveonline.com/regions/10000014/", "name": "Catch"}, {"href": "https://api-sisi.testeveonline.com/regions/10000051/", "name": "Cloud Ring"}, {"href": "https://api-sisi.testeveonline.com/regions/10000053/", "name": "Cobalt Edge"}, {"href": "https://api-sisi.testeveonline.com/regions/10000012/", "name": "Curse"}, {"href": "https://api-sisi.testeveonline.com/regions/11000016/", "name": "D-R00016"}, {"href": "https://api-sisi.testeveonline.com/regions/11000017/", "name": "D-R00017"}, {"href": "https://api-sisi.testeveonline.com/regions/11000018/", "name": "D-R00018"}, {"href": "https://api-sisi.testeveonline.com/regions/11000019/", "name": "D-R00019"}, {"href": "https://api-sisi.testeveonline.com/regions/11000020/", "name": "D-R00020"}, {"href": "https://api-sisi.testeveonline.com/regions/11000021/", "name": "D-R00021"}, {"href": "https://api-sisi.testeveonline.com/regions/11000022/", "name": "D-R00022"}, {"href": "https://api-sisi.testeveonline.com/regions/11000023/", "name": "D-R00023"}, {"href": "https://api-sisi.testeveonline.com/regions/10000035/", "name": "Deklein"}, {"href": "https://api-sisi.testeveonline.com/regions/10000060/", "name": "Delve"}, {"href": "https://api-sisi.testeveonline.com/regions/10000001/", "name": "Derelik"}, {"href": "https://api-sisi.testeveonline.com/regions/10000005/", "name": "Detorid"}, {"href": "https://api-sisi.testeveonline.com/regions/10000036/", "name": "Devoid"}, {"href": "https://api-sisi.testeveonline.com/regions/10000043/", "name": "Domain"}, {"href": "https://api-sisi.testeveonline.com/regions/11000024/", "name": "E-R00024"}, {"href": "https://api-sisi.testeveonline.com/regions/11000025/", "name": "E-R00025"}, {"href": "https://api-sisi.testeveonline.com/regions/11000026/", "name": "E-R00026"}, {"href": "https://api-sisi.testeveonline.com/regions/11000027/", "name": "E-R00027"}, {"href": "https://api-sisi.testeveonline.com/regions/11000028/", "name": "E-R00028"}, {"href": "https://api-sisi.testeveonline.com/regions/11000029/", "name": "E-R00029"}, {"href": "https://api-sisi.testeveonline.com/regions/10000039/", "name": "Esoteria"}, {"href": "https://api-sisi.testeveonline.com/regions/10000064/", "name": "Essence"}, {"href": "https://api-sisi.testeveonline.com/regions/10000027/", "name": "Etherium Reach"}, {"href": "https://api-sisi.testeveonline.com/regions/10000037/", "name": "Everyshore"}, {"href": "https://api-sisi.testeveonline.com/regions/11000030/", "name": "F-R00030"}, {"href": "https://api-sisi.testeveonline.com/regions/10000046/", "name": "Fade"}, {"href": "https://api-sisi.testeveonline.com/regions/10000056/", "name": "Feythabolis"}, {"href": "https://api-sisi.testeveonline.com/regions/10000058/", "name": "Fountain"}, {"href": "https://api-sisi.testeveonline.com/regions/11000031/", "name": "G-R00031"}, {"href": "https://api-sisi.testeveonline.com/regions/10000029/", "name": "Geminate"}, {"href": "https://api-sisi.testeveonline.com/regions/10000067/", "name": "Genesis"}, {"href": "https://api-sisi.testeveonline.com/regions/10000011/", "name": "Great Wildlands"}, {"href": "https://api-sisi.testeveonline.com/regions/11000032/", "name": "H-R00032"}, {"href": "https://api-sisi.testeveonline.com/regions/10000030/", "name": "Heimatar"}, {"href": "https://api-sisi.testeveonline.com/regions/10000025/", "name": "Immensea"}, {"href": "https://api-sisi.testeveonline.com/regions/10000031/", "name": "Impass"}, {"href": "https://api-sisi.testeveonline.com/regions/10000009/", "name": "Insmother"}, {"href": "https://api-sisi.testeveonline.com/regions/10000017/", "name": "J7HZ-F"}, {"href": "https://api-sisi.testeveonline.com/regions/10000052/", "name": "Kador"}, {"href": "https://api-sisi.testeveonline.com/regions/10000049/", "name": "Khanid"}, {"href": "https://api-sisi.testeveonline.com/regions/10000065/", "name": "Kor-Azor"}, {"href": "https://api-sisi.testeveonline.com/regions/10000016/", "name": "Lonetrek"}, {"href": "https://api-sisi.testeveonline.com/regions/10000013/", "name": "Malpais"}, {"href": "https://api-sisi.testeveonline.com/regions/10000042/", "name": "Metropolis"}, {"href": "https://api-sisi.testeveonline.com/regions/10000028/", "name": "Molden Heath"}, {"href": "https://api-sisi.testeveonline.com/regions/10000040/", "name": "Oasa"}, {"href": "https://api-sisi.testeveonline.com/regions/10000062/", "name": "Omist"}, {"href": "https://api-sisi.testeveonline.com/regions/10000021/", "name": "Outer Passage"}, {"href": "https://api-sisi.testeveonline.com/regions/10000057/", "name": "Outer Ring"}, {"href": "https://api-sisi.testeveonline.com/regions/10000059/", "name": "Paragon Soul"}, {"href": "https://api-sisi.testeveonline.com/regions/10000063/", "name": "Period Basis"}, {"href": "https://api-sisi.testeveonline.com/regions/10000066/", "name": "Perrigen Falls"}, {"href": "https://api-sisi.testeveonline.com/regions/10000048/", "name": "Placid"}, {"href": "https://api-sisi.testeveonline.com/regions/10000047/", "name": "Providence"}, {"href": "https://api-sisi.testeveonline.com/regions/10000023/", "name": "Pure Blind"}, {"href": "https://api-sisi.testeveonline.com/regions/10000050/", "name": "Querious"}, {"href": "https://api-sisi.testeveonline.com/regions/10000008/", "name": "Scalding Pass"}, {"href": "https://api-sisi.testeveonline.com/regions/10000032/", "name": "Sinq Laison"}, {"href": "https://api-sisi.testeveonline.com/regions/10000044/", "name": "Solitude"}, {"href": "https://api-sisi.testeveonline.com/regions/10000022/", "name": "Stain"}, {"href": "https://api-sisi.testeveonline.com/regions/10000041/", "name": "Syndicate"}, {"href": "https://api-sisi.testeveonline.com/regions/10000020/", "name": "Tash-Murkon"}, {"href": "https://api-sisi.testeveonline.com/regions/10000045/", "name": "Tenal"}, {"href": "https://api-sisi.testeveonline.com/regions/10000061/", "name": "Tenerifis"}, {"href": "https://api-sisi.testeveonline.com/regions/10000038/", "name": "The Bleak Lands"}, {"href": "https://api-sisi.testeveonline.com/regions/10000033/", "name": "The Citadel"}, {"href": "https://api-sisi.testeveonline.com/regions/10000002/", "name": "The Forge"}, {"href": "https://api-sisi.testeveonline.com/regions/10000034/", "name": "The Kalevala Expanse"}, {"href": "https://api-sisi.testeveonline.com/regions/10000018/", "name": "The Spire"}, {"href": "https://api-sisi.testeveonline.com/regions/10000010/", "name": "Tribute"}, {"href": "https://api-sisi.testeveonline.com/regions/10000004/", "name": "UUA-F4"}, {"href": "https://api-sisi.testeveonline.com/regions/10000003/", "name": "Vale of the Silent"}, {"href": "https://api-sisi.testeveonline.com/regions/10000015/", "name": "Venal"}, {"href": "https://api-sisi.testeveonline.com/regions/10000068/", "name": "Verge Vendor"}, {"href": "https://api-sisi.testeveonline.com/regions/10000006/", "name": "Wicked Creek"}], "pageCount": 1, "pageCount_str": "1", "totalCount": 99}`)

func contains(r []*Region, name string) *Region {
	for _, a := range r {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func TestUnpackRegions(t *testing.T) {
	regions := newRegions()
	page, err := unpackPage(regionsData)
	if err != nil {
		t.Error(err)
	}

	if page.hasNext {
		t.Error("Regions has a next page when it shouldn't")
	}

	regions, err = unpackRegions(regions, page)

	if err != nil {
		t.Error(err)
	}
	c := contains(regions.AllRegions, "The Forge")
	if c == nil {
		t.Error("Doesn't contain The Forge")
	}

	if c.Id != 10000002 {
		t.Error("Didn't decode the right region ID: %s", c)
	}

	if len(regions.AllRegions) != 99 {
		t.Error("Not enough regions")
	}
}

func TestUnpackRoot(t *testing.T) {
	root, err := unpackRoot(rootData)

	if err != nil {
		t.Error(err)
	}

	_, ok := root.Resources["regions"]
	if !ok {
		t.Error("Can't find a key for regions")
	}

}
