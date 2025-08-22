package shopee

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/donghui12/shopee_tool_base/pkg/logger"
	"github.com/donghui12/shopee_tool_base/pkg/pool"

	"go.uber.org/zap"
)

// func TestLogin(t *testing.T) {
//     client, config := setupTestClient(t)
// 	fmt.Println(config)

//     tests := []struct {
//         name     string
//         phone    string
//         password string
//         wantErr  bool
//     }{
//         {
//             name:     "正常登录",
//             phone:    config.Phone,
//             password: config.Password,
//             wantErr:  false,
//         },
//     }

//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             err := client.Login(tt.phone, tt.password)
//             if (err != nil) != tt.wantErr {
//                 t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
//             }
//         })
//     }
// }

func TestGetProductList(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()

	cookies := "SPC_CNSC_SESSION=5989a862fc6c74be0c3108cee22d607d_2_2306869_gS2ywaDg6rbNaKni2PWR48ZTWpiNrfufk1oR9IAWyJ3GRGPF9dbiyCateEAxfLytmQ2KDvk2VFX8Fo8uIE6GGY61c8YlbsMHYjoazdZjg0zO5bNyO7ai2x5KRq/rhMhrbobAAteNKeNrMOrMM1XpY8dYKMKXvuXtP7D7kn6kztlQ6abKrpxlxX0RvEETUXGHlgzxzlzOIwIiDSuFPrkIo1umfD6iDyHwQQBTNi4tlmlE=; CNSC_SSO=NHNyRHNXM2tIektZbE1nZAwSNTj06ceSjLToKi/revWm24XaYjSid8VWQQ07aFu6; SPC_SC_OFFLINE_TOKEN=eyJkYXRhIjoibEkvS2xGcndQNXY2dW8vRmVnbjJwbllCbDhGc1RUUCtGZFkxeFZ2REZDMGFac3lBa1RjdHgybzZ0K2Fsd1pEV0dzMWlDR1pDb1UxQUtVbW9ZOS84S1ZZU1FIN1oydUxHd2xFU2V6MnJNamlNRitieEtnRGYwUHR3clVmaXBuZkZYcFZmT1lEclc1VVFsakdJQTh2L0xSNlhtcmpqKzNObVpCYVVoMFZWNGd4ZEk0TjhCSllDbHBvUGhaemFFVlFLUzNVUnFrNC8vaHUxUHZFa05WN1o0Zz09IiwiaXYiOiJZV0JNYm13ZXRvK2loVGJabXZ2dEF3PT0iLCJzaWduIjoieGhOdCtjRk5Oa1NqZExudFB4Ykx4SkpLMnZQbG82WmZWQ1BnakFhcFVQRVhxeTQ5NXRGc2N5TWF2cUtYWThPQUp4di9yOTJUOVpGNTF1U2U3Z0ZwRFE9PSJ9; SPC_SI=dkNKaAAAAAA0TEpQZmpYaQHBKAAAAAAAWkt4b1paSjY=; SPC_SEC_SI=v1-UEFDRnhFQ3dVdURzeXlLZJSm26JICX5xJuuhTn8loWkP7JGQ0RBIzpzzd84CBJ8vNy1a0BFCorxFDLzGFIgqQBYWp8a9CdlWv5qUruFK6i0=;"
	shopID := "1313851163"
	region := "sg"

	productIDs, err := client.GetProductListWithDayToShip(cookies, shopID, region, "10", 10)
	if err != nil {
		t.Errorf("GetProductList() error = %v", err)
	}
	t.Logf("productIDs: %v", productIDs)
	t.Logf("len(productIDs): %d", len(productIDs))
}

func TestGetMerchantShopList(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	cookies := "SPC_CNSC_SESSION=f42893f98d19410156a15246b727dd03_2_2339367_g5Nw+/j4z0csq0Ef5PVP7/6DoKL56Fsmfqh/4WyMUak3RNv64Zb9qz5+dZhSyOqsQpUGeKLZBB5EmGrX5I1Y9GFEZcbB6kErRcf3oCspw12zftQzrlZY2ycBFJtCHJVUueUTAf0HNO4S35srfoIrcgWPgYodg4BZEGaISuAt0WorjxhzcAdVTDBUKt24kJsIaFiQimKVSJowwF5wUIWckls9kgTVwsavI5oqENE9TAjk=;"
	merchantShopList, err := client.GetMerchantShopList(cookies)
	if err != nil {
		t.Errorf("GetMerchantShopList() error = %v", err)
	}
	t.Logf("merchantShopList: %v", merchantShopList)
}

func TestUpdateProductInfo(t *testing.T) {
	client := GetShopeeClient()
	cookies := "SPC_CNSC_SESSION=c9ad3caf0d1d2d15d25d6e752a6c5723_2_2375038;"
	productID := int64(28760843741)
	day := 10
	shopID := "1350463893"
	region := "my"
	req := UpdateProductInfoReq{
		productID,
		day,
		cookies,
		shopID,
		region,
		ProductStatusInfo{},
	}
	err := client.UpdateProductInfo(req)
	if err != nil {
		t.Errorf("UpdateProductInfo() error = %v", err)
	}
}

func TestClient_GetProductListWithAreaTw(t *testing.T) {
	InitShopeeClient()
	client := GetTwShopeeClient()
	shopId := "15697232"
	accessToken := "785254686e655474424a5147464a6e4d"
	productIds, err := client.GetProductListWithAreaTw(accessToken, shopId)
	if err != nil {
		fmt.Printf("GetProductListWithAreaTw err: %s", err.Error())
		return
	}
	fmt.Println("this is productIds", productIds)
}

func TestClient_GetAccessTokenWithAreaTw(t *testing.T) {
	InitShopeeClient()
	client := GetTwShopeeClient()
	shopId := "1030427218"
	// code := "72544853557973664a4a586351685466"
	elderRefreshToken := "79437878736670584572554a73656778"
	accessToken, refreshToken, expireTimeFormatted, err := client.GetAccessTokenWithAreaTw(shopId, "", elderRefreshToken)
	if err != nil {
		fmt.Printf("GetProductListWithAreaTw err: %s", err.Error())
		return
	}
	fmt.Printf("this is accessToken: %s, refreshToken: %s", accessToken, refreshToken)
	fmt.Printf("this is expireTimeFormatted: %s", expireTimeFormatted)
}

func TestClient_UpdateProductInfoWithAreaTw(t *testing.T) {
	InitShopeeClient()
	client := GetTwShopeeClient()
	shopId := "1030427218"
	accessToken := "5a50436e764278756f6765644650614f"
	itemId := int64(28077511261)
	item := UpdateProductInfoWithAreaTwItem{
		DaysToShip: 2,
		IsPreOrder: false,
	}
	err := client.UpdateProductInfoWithAreaTw(accessToken, shopId, itemId, item)
	if err != nil {
		fmt.Printf("UpdateProductInfoWithAreaTw err: %s", err.Error())
		return
	}
}

func TestClient_GetProductInfoListWithAreaTw(t *testing.T) {
	InitShopeeClient()
	client := GetTwShopeeClient()
	shopId := "1030427218"
	accessToken := "5a50436e764278756f6765644650614f"
	shopInfoList := []int64{29076272050, 28976318873, 28877958509, 28777756753, 28776244241, 28775836504, 28527864423, 28476014973, 28427131288, 28126462353, 28125950071, 27926965997, 27926449591, 27627958449, 27627752969, 27627136075, 27427864049, 27426008882, 27425829239, 27377833725, 27377745494, 27277885959, 27126691930, 27126324359, 27027837691, 26977516589, 26975817418, 26926687443, 26877230704, 26727138958, 26577227137, 26425834082, 26076341896, 25541296785, 24591257650, 29927408216, 29826457164, 29826028466, 29777753402, 29675832473, 29627431951, 29626696467, 29427864462, 29325960336, 29077954483, 29075943590, 29027837882, 28775538512, 28727416585, 28426449238}
	currentProductItemList, err := client.GetProductBaseInfoWithAreaTw(accessToken, shopId, shopInfoList)
	if err != nil {
		logger.Error("获取当前商品基础信息失败", zap.Any("商品列表", shopInfoList))
	}
	fmt.Print("this is response, ", currentProductItemList)

}

func TestClient_GetInactiveProducts(t *testing.T) {
	InitShopeeClient()
	// 初始化线程池
	pool.InitWorkerPool()
	defer pool.GetWorkerPool().Release()
	client := GetShopeeClient()
	shopId := "1319428361"
	accessToken := "SPC_CNSC_SESSION=67651aad11a61bbf0b15c7d1bfca2aa9_2_2316152;"
	currentProductItemList, err := client.GetInactiveProducts(accessToken, shopId, "sg", 50)
	if err != nil {
		logger.Error("获取当前商品基础信息失败")
	}
	fmt.Print("this is response, ", currentProductItemList)
}

func TestClient_DeleteProducts(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1319428361"
	cookies := "SPC_CNSC_SESSION=67651aad11a61bbf0b15c7d1bfca2aa9_2_2316152;"
	var ProductIdList []int64
	ProductIdList = append(ProductIdList, 28365112846)
	success_count, _ := client.DeleteProducts(shopId, cookies, "sg", ProductIdList)
	fmt.Printf("成功：%d\n", success_count)
}

func TestClient_GetDiscountList(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1332278997"
	cookieStr := "SPC_CNSC_SESSION=f42893f98d19410156a15246b727dd03_2_2339367_g5Nw+/j4z0csq0Ef5PVP7/6DoKL56Fsmfqh/4WyMUak3RNv64Zb9qz5+dZhSyOqsQpUGeKLZBB5EmGrX5I1Y9GFEZcbB6kErRcf3oCspw12zftQzrlZY2ycBFJtCHJVUueUTAf0HNO4S35srfoIrcgWPgYodg4BZEGaISuAt0WorjxhzcAdVTDBUKt24kJsIaFiQimKVSJowwF5wUIWckls9kgTVwsavI5oqENE9TAjk=;"

	discounts, err := client.GetDiscountList(cookieStr, shopId, "sg")
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range discounts {
		endTime := d.SellerDiscount.EndTime
		t := time.Unix(endTime, 0).Local().Format("2006-01-02 15:04:05")

		fmt.Printf("折扣Id: %d, 折扣名称: %s, 折扣状态: %d, 结束时间: %s, 商品数量: %d\n",
			d.SellerDiscount.DiscountID,
			d.SellerDiscount.Name,
			d.SellerDiscount.TimeStatus,
			t,
			d.SellerDiscount.ItemPreview.ItemCount)
	}
	fmt.Printf("总计:%d\n", len(discounts))
}

func TestClient_DeleteDiscounts(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1516433576"
	cookieStr := "SPC_CNSC_SESSION=d70d1113b908bcd83652d52a08ac6d3b_2_3104665_gccNRk8lTl4PQbrZbJgwafj2Vdnyp5kNlCnEuk3yFQCp5mAxXMgBG+PJQ9sxc2F4eRkkX0eCMFU9YrVwI3Xb4TWavqo0JXY8RfF/GxYLmWUi2O2A8LmBjgUcprztFmdVfJvVYEAyrxzV06EFBTZbQb1yDPdZJQr39b2R9KvK0WdPQi4BxQL+EAE9bQX3BgfdR1/MWaZ/zoKHFPoAy9nlRLVUjjGJC4h6XD/BZkBdQ28c=;"
	discountList := []int64{}
	discountList = append(discountList, 641044027424768)
	successCount, err := client.DeleteDiscounts(cookieStr, shopId, "sg", discountList, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("successCount:", successCount)
}

func TestClient_GetDiscountItemList(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1586711525"
	cookieStr := "SPC_CNSC_SESSION=eff1ee3c3314f86f9cc3119b3decd2a6_2_4606533_g6XFZrMuK/LhI+7uoeFcXpsy79055JdECTu9HfOE+3Z6IsQPirNdkgC/SMXIKNLJBgzWJq4w71R+YuJc4nJS0Hgu7nptCtVACnRZ4M/gbKPR6F/gLZx35cgRmDWUH+1ICh7YRMDuc3tkbO+zsPODCcPLRiGEPJfo9/HJ3Td/hr/9MtIVm860y5JyZsncDAW53SvlaQH+ULRcOOvA3OizePjfxXvzmlG97BGcbBrpVQHk=; CNSC_SSO=aHNwcjNMQzlkT3FjenNTMRarJeKt4zkvC8Lxi+zHHF98i0PiDZ9vawyZYLoigiWD; SPC_SC_OFFLINE_TOKEN=eyJkYXRhIjoicExiLzJDU1FiZCtMU1oyZEVpcXhNSjhmZUdOM3JMSEhMOWh4ZDl1SGRqTDkxTjYwd0swTUFXUHlpSXNJeGFINHVSRDV6aVVvQ0UvdUVIYTJXZDJ4TGpwT3dwQ28wSkJuMnhRMnlHSHNEQTdqajNHWDR6M2RrR0RUMTBHWWdjVkZvZ3lkNlF2ZFh1Z3VlT2o0WjdNMHA2cnhRbEYwcHdnSUdEZHVDb2drZXZuVHIxd3VlbENFMTFhcEVtQVp3UCszTmNHYzNBRHI2R2FJeHJyYTkzUjBWQT09IiwiaXYiOiJXanF2YnlXMGl4WFlHV3g2TkFYKzBRPT0iLCJzaWduIjoiaFJFeWM4bXg2OVBjMzcyc2JET3hUUXRZS3NHSFpnWHB6OGZSL3FwS0o5RVdINE1CZ1BnODRVUEZXNHMyZWlUK1lIcXdSRHFoUDM3RFhIUHNNaGQ0Q2c9PSJ9; SPC_SI=H4uHaAAAAAB3Z29iUlQyduA+FQAAAAAAdXRJMU5valc=; SPC_SEC_SI=v1-TE1PWUR5bmRKUm15NXM3YaOZ0cd7FwHMOD8gVAKxLIYjSX3ATndyhvf63f14dthLErMHjxUGz6d6THBnu2wYXTyT2y2lGj0Vnjc5teQYSzo=;"
	// var discountId int64
	discountId := int64(707239305150464)
	discountItemList, err := client.GetDiscountItem(cookieStr, shopId, "sg", discountId)
	if err != nil {
		log.Fatal(err)
	}
	itemIdMap := make(map[int64]int)
	for _, d := range discountItemList {
		itemIdMap[d.ItemID] += 1
		// fmt.Printf("ItemId: %d, ModelId: %d, 价格: %d, 优惠: %d, 状态: %d, UserItemLimit: %d\n",
		// 	d.ItemID,
		// 	d.ModelID,
		// 	d.PromotionPrice,
		// 	d.PromotionStock,
		// 	d.Status,
		// 	d.UserItemLimit,
		// )
	}
	fmt.Printf("this is itemId:%d\n", len(itemIdMap))
	fmt.Printf("总计:%d\n", len(discountItemList))
}

func TestClient_UpdateDiscountItem(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1586711506"
	cookieStr := "SPC_CNSC_SESSION=eff1ee3c3314f86f9cc3119b3decd2a6_2_4606533_g6XFZrMuK/LhI+7uoeFcXpsy79055JdECTu9HfOE+3Z6IsQPirNdkgC/SMXIKNLJBgzWJq4w71R+YuJc4nJS0Hgu7nptCtVACnRZ4M/gbKPR6F/gLZx35cgRmDWUH+1ICh7YRMDuc3tkbO+zsPODCcPLRiGEPJfo9/HJ3Td/hr/9MtIVm860y5JyZsncDAW53SvlaQH+ULRcOOvA3OizePjfxXvzmlG97BGcbBrpVQHk=;"

	discountId := int64(708793210257408)
	discountItemList, err := client.GetDiscountItem(cookieStr, shopId, "sg", discountId)
	if err != nil {
		log.Fatal(err)
	}
	req := UpdateDiscountItemRequest{
		PromotionId:       708795332575232,
		DiscountModelList: discountItemList,
	}
	successCount, err := client.UpdateDiscountItem(cookieStr, shopId, "sg", req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("总计:%d\n", successCount)
}

func TestClient_CopyCreateDiscounts(t *testing.T) {
	InitShopeeClient()
	client := GetShopeeClient()
	shopId := "1586711506"
	cookieStr := "SPC_CNSC_SESSION=eff1ee3c3314f86f9cc3119b3decd2a6_2_4606533_g6XFZrMuK/LhI+7uoeFcXpsy79055JdECTu9HfOE+3Z6IsQPirNdkgC/SMXIKNLJBgzWJq4w71R+YuJc4nJS0Hgu7nptCtVACnRZ4M/gbKPR6F/gLZx35cgRmDWUH+1ICh7YRMDuc3tkbO+zsPODCcPLRiGEPJfo9/HJ3Td/hr/9MtIVm860y5JyZsncDAW53SvlaQH+ULRcOOvA3OizePjfxXvzmlG97BGcbBrpVQHk=; CNSC_SSO=aHNwcjNMQzlkT3FjenNTMRarJeKt4zkvC8Lxi+zHHF98i0PiDZ9vawyZYLoigiWD; SPC_SC_OFFLINE_TOKEN=eyJkYXRhIjoicExiLzJDU1FiZCtMU1oyZEVpcXhNSjhmZUdOM3JMSEhMOWh4ZDl1SGRqTDkxTjYwd0swTUFXUHlpSXNJeGFINHVSRDV6aVVvQ0UvdUVIYTJXZDJ4TGpwT3dwQ28wSkJuMnhRMnlHSHNEQTdqajNHWDR6M2RrR0RUMTBHWWdjVkZvZ3lkNlF2ZFh1Z3VlT2o0WjdNMHA2cnhRbEYwcHdnSUdEZHVDb2drZXZuVHIxd3VlbENFMTFhcEVtQVp3UCszTmNHYzNBRHI2R2FJeHJyYTkzUjBWQT09IiwiaXYiOiJXanF2YnlXMGl4WFlHV3g2TkFYKzBRPT0iLCJzaWduIjoiaFJFeWM4bXg2OVBjMzcyc2JET3hUUXRZS3NHSFpnWHB6OGZSL3FwS0o5RVdINE1CZ1BnODRVUEZXNHMyZWlUK1lIcXdSRHFoUDM3RFhIUHNNaGQ0Q2c9PSJ9; SPC_SI=H4uHaAAAAAB3Z29iUlQyduA+FQAAAAAAdXRJMU5valc=; SPC_SEC_SI=v1-TE1PWUR5bmRKUm15NXM3YaOZ0cd7FwHMOD8gVAKxLIYjSX3ATndyhvf63f14dthLErMHjxUGz6d6THBnu2wYXTyT2y2lGj0Vnjc5teQYSzo=;"
	discountList := []Discount{}
	discountList = append(discountList, Discount{
		DiscountType: 1,
		SellerDiscount: SellerDiscount{
			DiscountID: 704330479190016,
			Name:       "7.23",
			TimeStatus: 3,
			StartTime:  1753268400,
			EndTime:    1753607624,
			ItemPreview: ItemPreview{
				ItemCount: 198,
				Images: []string{
					"sg-11134201-7rdy2-mcj4lrjw14mr9f",
					"sg-11134201-7rdw1-mcdl1c39iul700",
					"sg-11134201-7rdxs-mch7c31n9h04a9",
					"sg-11134201-7rdwn-mcj4gsl601ye93",
					"sg-11134201-7rdvi-mcj4nlxeuwjkd0",
				},
			},
		},
	})
	successCount, err := client.CopyCreateDiscount(cookieStr, shopId, "sg", discountList[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("successCount:", successCount)
}
