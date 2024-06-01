package local_cache

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

//var lc = NewLocalCache(100, 1000, 2*time.Second)

func BenchmarkLocalCache_Set(b *testing.B) {
	lc := NewLocalCache(10000, 1000, 60*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
	}
	b.StopTimer()
}

func BenchmarkLocalCache_LruSet(b *testing.B) {
	lc := NewLocalCache(1, 15, 5*time.Second)
	for i := 0; i < 15; i++ {
		lc.Set(fmt.Sprintf("%v", i), i, NoExpiration)
	}
	//for i := 0; i < 10; i++ {
	//	lc.Set(fmt.Sprintf("%v", i), i, NoExpiration)
	//}
	lc.Get("1")
	lc.bucketsDta[0].ListLRUCache()
}

func BenchmarkLocalCache_Get(b *testing.B) {

	lc := NewLocalCache(100, 100, 10*time.Second)
	for i := 0; i < b.N; i++ {
		lc.Get(GenerateRandomString(10))
	}
}

const data = "{\"code\":0,\"msg\":\"OK\",\"data\":{\"alive_conf\":{\"alive_mode\":0,\"alive_tab\":0,\"alive_type_state\":1,\"anti_screen_capture\":1,\"app_anti_screen_capture\":1,\"app_anti_screen_jump\":0,\"authentic_state\":1,\"can_record\":0,\"complete_time\":0,\"cutting_type\":0,\"eclockwm_anti_screen_jump\":0,\"elearn_anti_screen_capture\":1,\"elearn_anti_screen_jump\":0,\"elive_anti_screen_jump\":0,\"elive_anti_screen_jump_url\":\"\",\"esnswm_anti_screen_jump\":0,\"esnswmd_anti_screen_jump\":0,\"forbid_record\":1,\"forbid_talk\":0,\"guide_binding_phone_switch\":1,\"h5_anti_screen_jump\":0,\"h5_anti_screen_jump_url\":\"\",\"h5_url\":\"https://appAKLWLitn7978.h5.xiaoeknow.com\",\"has_invite\":1,\"has_reward\":1,\"if_push\":1,\"is_audit_first_on\":0,\"is_can_exceptional\":0,\"is_card_on\":0,\"is_coupon_on\":0,\"is_heat_on\":0,\"is_hor_float_card_on\":1,\"is_invite_data_on\":0,\"is_invite_on\":1,\"is_jump_full_screen\":1,\"is_lookback\":1,\"is_message_on\":0,\"is_online_on\":1,\"is_open_complete_time\":0,\"is_open_ego_mode\":0,\"is_open_promoter\":0,\"is_open_qus\":1,\"is_open_share_reward\":0,\"is_open_task\":0,\"is_open_vote\":1,\"is_picture_on\":0,\"is_point_red_packet_on\":0,\"is_privacy_protection\":1,\"is_prize_on\":0,\"is_red_packet_on\":0,\"is_redirect_index\":0,\"is_round_table_on\":1,\"is_set_preview\":0,\"is_show_marquee\":1,\"is_show_reward\":0,\"is_show_reward_on\":0,\"is_show_view_count\":0,\"is_sign_in_on\":1,\"is_takegoods\":0,\"is_thumb_on\":1,\"lookback_time\":{\"expire\":-1,\"expire_type\":1},\"msg_bubble_type\":1,\"only_h5_play\":0,\"open_pc_network_school\":1,\"pc_network_school_index_url\":\"pc.inside.xiaoecloud.com\",\"privacy_protection_live\":2,\"red_packet_switch\":0,\"relate_sell_info\":1,\"reward_switch\":1,\"share_file_switch\":1,\"show_on_wall\":0,\"show_on_wall_switch\":1,\"version_type\":4,\"video_player_type\":1,\"view_in_mini_program\":1,\"view_stop_switch\":1,\"warm_up\":1,\"wm_anti_screen_jump\":0,\"wx_app_avatar\":\"https://wechatapppro-1252524126.file.myqcloud.com/appAKLWLitn7978/image/kcbd0li707pmulc8q5ei.jpg\",\"wx_app_name\":\"现网蓝悦2号【7978勿动】\"},\"alive_info\":{\"alive_id\":\"l_6618e9c3e4b023c0a96a3285\",\"alive_img_url\":\"https://wechatapppro-1252524126.file.myqcloud.com/appAKLWLitn7978/image/kquiu79x0pa8.jpg\",\"alive_room_url\":\"https://appAKLWLitn7978.h5.xiaoeknow.com/v2/course/alive/l_6618e9c3e4b023c0a96a3285?app_id=appAKLWLitn7978\\u0026pro_id=\\u0026type=2\",\"alive_state\":3,\"alive_type\":2,\"app_id\":\"appAKLWLitn7978\",\"can_select\":1,\"checktimestamp\":1716901815,\"comment_count\":0,\"cover_img_url\":\"https://commonresource-1252524126.cdn.xiaoeknow.com/image/l6nfw9120t1u.png\",\"create_mode\":0,\"descrb\":\"\",\"img_url\":\"https://commonresource-1252524126.cdn.xiaoeknow.com/image/l6nfw9120t1u.png\",\"img_url_compressed\":\"https://commonresource-1252524126.cdn.xiaoeknow.com/image/l6nfw9120t1u.png\",\"manual_stop_at\":null,\"old_live_room_url\":\"/content_page/eyJ0eXBlIjoxMiwicmVzb3VyY2VfdHlwZSI6NCwicmVzb3VyY2VfaWQiOiJsXzY2MThlOWMzZTRiMDIzYzBhOTZhMzI4NSIsInByb2R1Y3RfaWQiOiIiLCJwYXltZW50X3R5cGUiOjEsImNoYW5uZWxfaWQiOiIiLCJhcHBfaWQiOiJhcHBBS0xXTGl0bjc5NzgiLCJzb3VyY2UiOiIiLCJzY2VuZSI6IiIsImNvbnRlbnRfYXBwX2lkIjoiIiwiZXh0cmFfZGF0YSI6MCwid2ViX2FsaXZlIjowLCJ0b2tlbiI6IiJ9\",\"org_content\":\"\",\"param_str\":\"eyJwYXltZW50X3R5cGUiOjIsInByb2R1Y3RfaWQiOiIiLCJyZXNvdXJjZV9pZCI6ImxfNjYxOGU5YzNlNGIwMjNjMGE5NmEzMjg1IiwicmVzb3VyY2VfdHlwZSI6NH0\",\"product_id\":\"\",\"product_name\":\"\",\"push_ahead\":\"-1\",\"push_state\":2,\"push_url\":\"\",\"pushzb_start_at\":\"2024-04-12 15:59:35\",\"pushzb_stop_at\":\"2024-04-12 16:59:35\",\"record_push_end_time\":\"2024-04-12 15:58:35\",\"remainder_time\":0,\"resource_type\":4,\"room_id\":\"XET#898b0d09d800ec9\",\"sell_mode\":0,\"summary\":\"快来观看我的直播～\",\"title\":\"CK的直播间\",\"user_title\":\"\",\"user_type\":0,\"view_count\":0,\"zb_countdown_time\":-3993100,\"zb_start_at\":1712908775,\"zb_stop_at\":1712912375},\"alive_play\":{\"pc_alive_video_url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.flv\",\"mini_alive_video_url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\",\"alive_video_url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\",\"new_alive_host\":\"live-flexible.xiaoeknow.com/live/\",\"alive_fast_webrtcurl\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP\",\"new_alive_video_url\":\"\",\"fast_alive_switch\":true,\"video_alive_use_cos\":false,\"alive_video_more_sharpness\":[{\"definition_name\":\"原画\",\"definition_p\":\"default\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\"},{\"definition_name\":\"超清\",\"definition_p\":\"sd\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\"},{\"definition_name\":\"高清\",\"definition_p\":\"hd\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"},{\"definition_name\":\"流畅\",\"definition_p\":\"fluent\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"}],\"pc_alive_video_more_sharpness\":[{\"definition_name\":\"原画\",\"definition_p\":\"default\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.flv\"},{\"definition_name\":\"超清\",\"definition_p\":\"sd\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.flv\"},{\"definition_name\":\"高清\",\"definition_p\":\"hd\",\"encrypt\":\"\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\"},{\"definition_name\":\"流畅\",\"definition_p\":\"fluent\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\"}],\"alive_fast_more_sharpness\":[{\"definition_name\":\"原画\",\"definition_p\":\"default\",\"encrypt\":\"\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP\"},{\"definition_name\":\"超清\",\"definition_p\":\"sd\",\"encrypt\":\"\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P\"},{\"definition_name\":\"高清\",\"definition_p\":\"hd\",\"encrypt\":\"\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P\"},{\"definition_name\":\"流畅\",\"definition_p\":\"fluent\",\"encrypt\":\"\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P\"}],\"recorded_use_pull_stream\":false,\"alive_video_backup_more_sharpness\":[[{\"definition_name\":\"原画\",\"definition_p\":\"default\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\"},{\"definition_name\":\"超清\",\"definition_p\":\"sd\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\"},{\"definition_name\":\"高清\",\"definition_p\":\"hd\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"},{\"definition_name\":\"流畅\",\"definition_p\":\"fluent\",\"encrypt\":\"\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"}],[{\"definition_name\":\"原画\",\"definition_p\":\"default\",\"encrypt\":\"\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\"},{\"definition_name\":\"超清\",\"definition_p\":\"sd\",\"encrypt\":\"\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\"},{\"definition_name\":\"高清\",\"definition_p\":\"hd\",\"encrypt\":\"\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"},{\"definition_name\":\"流畅\",\"definition_p\":\"fluent\",\"encrypt\":\"\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\"}]],\"lines\":[{\"line_name\":\"线路2\",\"default\":false,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\",\"default\":true,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"}]},{\"line_name\":\"线路3\",\"default\":false,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\",\"default\":true,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"}]},{\"line_name\":\"线路4\",\"default\":true,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP\",\"default\":true,\"cloud\":\"tx\",\"type\":\"fast\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P\",\"default\":false,\"cloud\":\"tx\",\"type\":\"fast\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P\",\"default\":false,\"cloud\":\"tx\",\"type\":\"fast\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"webrtc://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P\",\"default\":false,\"cloud\":\"tx\",\"type\":\"fast\"}]},{\"line_name\":\"线路5\",\"default\":false,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.m3u8\",\"default\":true,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.m3u8\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay-hw-hls.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.m3u8\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"}]}],\"elive_lines\":[{\"line_name\":\"线路1\",\"default\":false,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.flv\",\"default\":true,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.flv\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay-byte.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"byte\",\"type\":\"normal\"}]},{\"line_name\":\"线路2\",\"default\":false,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.flv\",\"default\":true,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.flv\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"tx\",\"type\":\"normal\"}]},{\"line_name\":\"线路3\",\"default\":true,\"line_sharpness\":[{\"name\":\"原画\",\"resolution\":\"origin\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP.flv\",\"default\":true,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"超清\",\"resolution\":\"sd\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_1080P.flv\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"高清\",\"resolution\":\"hd\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"},{\"name\":\"流畅\",\"resolution\":\"fluency\",\"url\":\"http://liveplay-hw-flv.xiaoeknow.com/live/5060_e_898b0d08c8003a3zP_720P.flv\",\"default\":false,\"cloud\":\"hw\",\"type\":\"normal\"}]}]},\"available_info\":{\"available\":true,\"available_product\":false,\"expire_at\":\"\",\"have_password\":0,\"is_public\":1,\"is_stop_sell\":0,\"is_try\":0,\"payment_type\":1,\"recycle_bin_state\":1},\"campro_report\":false,\"caption_define\":{\"audio_try_hint\":\"购买\",\"column_open\":\"订阅\",\"column_pay_hint\":\"开通会员\",\"column_title\":\"专栏\",\"home_tab_message\":\"消息\",\"home_title\":\"图文导航\",\"single_product_hint\":\"订阅专栏\"},\"e_course_data\":null,\"index_url\":\"https://appAKLWLitn7978.h5.xiaoeknow.com/homepage\",\"is_static_switch\":false,\"live_skin\":{\"h5_link\":\"\",\"key\":\"\",\"type\":0},\"payment_url\":\"\",\"share_info\":{\"share_info\":{\"is_share_free\":0,\"share_user_id\":\"\",\"num\":0,\"surplus_num\":0,\"share_resource\":0,\"has_share_resource\":null,\"product_info\":null},\"share_listen_info\":{\"is_share_listen\":false,\"is_show_share_count\":true}}},\"requestId\":\"2ef465e9619ac786\"}\n"

func BenchmarkAsynMemDel(b *testing.B) {

	printMem("begin")
	lc := NewLocalCache(100, 1000, 5*time.Second)
	printMem("new mem")

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			//time.Sleep(time.Second * time.Duration(rand.Intn(10)+1))
			for ii := 0; ii < 10000; ii++ {
				lc.Set(GenerateRandomString(10), data, time.Second*2)
			}
			wg.Done()
		}()
	}
	time.Sleep(time.Second * 10)
	printMem("after set")
	for {
		time.Sleep(time.Second * 5)
		lc.Get(GenerateRandomString(10))
		runtime.GC()
		printMem("after gc")

	}

}

func BenchmarkAsynDel(b *testing.B) {

	lc := NewLocalCache(1, 10, 2*time.Second)
	for i := 0; i < 100000; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			lc.Set(fmt.Sprintf("alive_id_%d", i), i, NoExpiration)
		}(i)

		go func(i int) {
			defer wg.Done()
			lc.Get(fmt.Sprintf("alive_id_%d", i))
		}(i)
	}
	wg.Wait()
	fmt.Println(lc.bucketsDta[0].items)
	lc.bucketsDta[0].ListLRUCache()
	time.Sleep(4 * time.Second)

	fmt.Println(lc.bucketsDta[0].items)
	lc.bucketsDta[0].ListLRUCache()
	fmt.Println("1")
}

func BenchmarkLruReadDel(b *testing.B) {
	lc := NewLocalCache(100, 1000, 5*time.Second)
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			lc.Set(fmt.Sprintf("alive_id_%d", i), i, NoExpiration)
		}(i)

		go func(i int) {
			defer wg.Done()
			lc.Get(fmt.Sprintf("alive_id_%d", i))
		}(i)
	}
	wg.Wait()
}

func BenchmarkLruLocalCache(b *testing.B) {

	lc := NewLocalCache(1, 10, 5*time.Second)
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 1)
		lc.Set(fmt.Sprintf("alive_id_xxx%v", i), i, NoExpiration)
	}
}

func printMem(str string) {

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println(fmt.Sprintf("\n %s MemAlloc: %dKB, TotalAlloc:%dKB, HeapAlloc:%dKB, HeapInuse:%dKB, StackInuse:%dKB, GC:%d次", str, mem.Alloc/1024, mem.TotalAlloc/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.StackInuse/1024, mem.NumGC))
}

func (b *bucket) ListLRUCache() {
	node := b.head.nextI
	for node != nil {
		fmt.Println(fmt.Sprintf("key: %s, value: %d", node.key, node.value))
		node = node.nextI
	}
	fmt.Println("xxxxxxxxxxxxxxxxx\\n")
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	// 使用加密安全的随机数生成器初始化种子
	rand.Seed(time.Now().UnixNano())
	// 使用加密安全的随机数填充切片
	rand.Read(b)
	// 将字节切片转换为base64字符串
	return base64.StdEncoding.EncodeToString(b)
}
