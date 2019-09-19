package conversion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fonzie1006/shortlink/pkg/gredis"
	"github.com/fonzie1006/shortlink/pkg/util"
	logging "github.com/sirupsen/logrus"
)

const (
	// redis incr value: 自增ID
	sid = "slink:sid"
	// slink:md5:<incrID>:<md5> value: shortlink
	lUrlMd5 = "slink:md5:%s"
	// slink:short:<shortlink> value: {"longlink":"","md5":"","incrid":"", "count":"", "shortlink":""}
	sUrl       = "slink:id:%d"
	urlMinSize = 14776334
	urlMaxLen  = 6
)

type ShortLinkConversion struct {
	LongLink  string `json:"longlink"`
	Shortlink string `json:"shortlink"`
	IncrID    int64  `json:"incrid"`
	// 长链接的md5
	Md5     string `json:"md5"`
	Count   int64  `json:"count"`
	Timeout int    `json:"timeout"`
}

// 长链接转换成短链接
func (s *ShortLinkConversion) GetShortLink() error {
	var err error
	if err = util.CheckUrl(s.LongLink); err != nil {
		return err
	}

	s.Md5 = util.EncodeMd5(s.LongLink)
	// 如果长链接已经存在从redis中直接获取
	if gredis.Exists(fmt.Sprintf(lUrlMd5, s.Md5)) {
		err = s.FromRedisGetFullInfo()
		if err != nil {
			return err
		}
	} else {
		err = s.GenerateShortLink()
		if err != nil {
			return err
		}
	}

	return nil
}

// 从redis获取长短链接的详细信息
func (s *ShortLinkConversion) FromRedisGetFullInfo() error {
	var err error
	// 通过长链接的md5获取短链接
	s.IncrID, err = gredis.GetInt(fmt.Sprintf(lUrlMd5, s.Md5))
	if err != nil {
		logging.Errorf("Get an error with a short link corresponding to an existing long link : %v", err)
		return err
	}

	fullInfoByte, err := gredis.Get(fmt.Sprintf(sUrl, s.IncrID))
	if err != nil {
		logging.Errorf("Get short link details error : %v", err)
		return err
	}

	err = json.Unmarshal(fullInfoByte, s)
	if err != nil {
		logging.Errorf("Parsing long links and short links complete information error: %v", err)
		return err
	}

	return nil
}

func (s *ShortLinkConversion) GenerateShortLink() error {
	flag := false
	n := 0

	for !flag {
		// 防止死循环
		if n > 5 {
			logging.Error("generate short link error: The number of cycles exceeded 5 times")
			return errors.New("generate short link error: The number of cycles exceeded 5 times")
		}

		incrID, err := gredis.Incr(sid)
		if err != nil {
			logging.Errorf("Failed to get incr from redis : %v", err)
			return err
		}

		if incrID < urlMinSize {
			err := gredis.Set(sid, urlMinSize, -1)
			if err != nil {
				logging.Errorf("reset redis incr id err : %v", err)
				return err
			}
			incrID = urlMinSize
			continue
		}

		s.Shortlink, err = util.Base62(incrID)
		if err != nil {
			logging.Errorf("get base62 err : %v", err)
			return err
		}

		if len(s.Shortlink) > urlMaxLen {
			err := gredis.Set(sid, urlMinSize, -1)
			if err != nil {
				logging.Errorf("Failed to get incr from redis : %v", err)
				return err
			}
			continue
		}

		s.IncrID = incrID

		// 保存之前先判断incrID是否已经存在如果已经存在则删除说明这个一轮大循环已经完成了,为了保证redis的key不会无限累加删除之前用过的短链接
		if gredis.Exists(fmt.Sprintf(sUrl, s.IncrID)) {
			tUrl := ShortLinkConversion{}
			tUrlByte, _ := gredis.Get(fmt.Sprintf(sUrl, s.IncrID))

			err = json.Unmarshal(tUrlByte, &tUrl)
			if err != nil {
				logging.Error(err)
				return err
			}

			_, err := gredis.Delete(fmt.Sprintf(lUrlMd5, tUrl.Md5))
			if err != nil {
				logging.Errorf("删除长链接Md5发生错误 : %v", err)
				return err
			}

		}

		// 保存短链接的详细信息
		err = gredis.Set(fmt.Sprintf(sUrl, s.IncrID), s, s.Timeout)
		if err != nil {
			logging.Errorf("Save short link details error : %v", err)
			return err
		}

		err = gredis.Set(fmt.Sprintf(lUrlMd5, s.Md5), s.IncrID, s.Timeout)
		if err != nil {
			logging.Errorf("Save long link md5 details error : %v", err)
			return err
		}

		flag = true
	}

	return nil
}

// 短链接转长链接
func (s *ShortLinkConversion) GetLongLink() error {
	sLinkId, err := util.Decode62(s.Shortlink)
	if err != nil {
		logging.Errorf("decode62 err : %v", err)
		return err
	}

	tUrlByte, err := gredis.Get(fmt.Sprintf(sUrl, sLinkId))
	if err != nil {
		logging.Errorf("Error getting short link details by id", err)
		return err
	}
	err = json.Unmarshal(tUrlByte, s)
	if err != nil {
		logging.Error(err)
		return err
	}

	return nil

}

// 记录短链接访问次数
func (s *ShortLinkConversion) CountShortLinkAccess() error {
	s.Count += 1
	err := gredis.Set(fmt.Sprintf(sUrl, s.IncrID), s, -1)
	if err != nil {
		logging.Errorf("记录短链接访问次数发生错误: %v", err)
		return err
	}
	return nil
}
