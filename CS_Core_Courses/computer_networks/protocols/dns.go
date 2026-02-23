package protocols

import (
	"fmt"
	"strings"
	"time"
)

// DNSRecordType DNS 记录类型
type DNSRecordType string

const (
	A     DNSRecordType = "A"     // IPv4 地址
	AAAA  DNSRecordType = "AAAA"  // IPv6 地址
	CNAME DNSRecordType = "CNAME" // 别名
	MX    DNSRecordType = "MX"    // 邮件服务器
	NS    DNSRecordType = "NS"    // 名称服务器
)

// DNSRecord DNS 记录
type DNSRecord struct {
	Name  string        // 域名
	Type  DNSRecordType // 记录类型
	Value string        // 记录值
	TTL   int           // 生存时间 (秒)
}

// String 格式化输出 DNS 记录
func (r *DNSRecord) String() string {
	return fmt.Sprintf("%-25s  TTL=%-6d  %-6s  %s", r.Name, r.TTL, r.Type, r.Value)
}

// DNSCache DNS 缓存
type DNSCache struct {
	Records   map[string]map[DNSRecordType]*DNSRecord // 域名 -> 类型 -> 记录
	Timestamp map[string]time.Time                    // 记录添加时间
}

// NewDNSCache 创建 DNS 缓存
func NewDNSCache() *DNSCache {
	return &DNSCache{
		Records:   make(map[string]map[DNSRecordType]*DNSRecord),
		Timestamp: make(map[string]time.Time),
	}
}

// Add 添加记录到缓存
func (c *DNSCache) Add(record *DNSRecord) {
	if c.Records[record.Name] == nil {
		c.Records[record.Name] = make(map[DNSRecordType]*DNSRecord)
	}
	c.Records[record.Name][record.Type] = record
	key := record.Name + string(record.Type)
	c.Timestamp[key] = time.Now()
}

// Lookup 查询缓存
func (c *DNSCache) Lookup(name string, recordType DNSRecordType) (*DNSRecord, bool) {
	if c.Records[name] == nil {
		return nil, false
	}

	record, exists := c.Records[name][recordType]
	if !exists {
		return nil, false
	}

	// 检查 TTL 是否过期
	key := name + string(recordType)
	if time.Since(c.Timestamp[key]).Seconds() > float64(record.TTL) {
		delete(c.Records[name], recordType)
		delete(c.Timestamp, key)
		return nil, false
	}

	return record, true
}

// CleanExpired 清除过期记录
func (c *DNSCache) CleanExpired() int {
	count := 0
	for name, typeMap := range c.Records {
		for recordType, record := range typeMap {
			key := name + string(recordType)
			if time.Since(c.Timestamp[key]).Seconds() > float64(record.TTL) {
				delete(typeMap, recordType)
				delete(c.Timestamp, key)
				count++
			}
		}
		if len(typeMap) == 0 {
			delete(c.Records, name)
		}
	}
	return count
}

// PrintCache 打印缓存
func (c *DNSCache) PrintCache() {
	fmt.Println("\n【DNS 缓存】")
	if len(c.Records) == 0 {
		fmt.Println("  (空)")
		return
	}

	fmt.Println("┌───────────────────────────────────────────────────────────────┐")
	fmt.Println("│ 域名                      TTL     类型    值                 │")
	fmt.Println("├───────────────────────────────────────────────────────────────┤")
	for _, typeMap := range c.Records {
		for _, record := range typeMap {
			fmt.Printf("│ %s │\n", record)
		}
	}
	fmt.Println("└───────────────────────────────────────────────────────────────┘")
}

// DNSServer DNS 服务器
type DNSServer struct {
	Name    string                                      // 服务器名称
	Records map[string]map[DNSRecordType]*DNSRecord    // DNS 记录数据库
	Parent  *DNSServer                                  // 父级服务器 (用于迭代查询)
}

// NewDNSServer 创建 DNS 服务器
func NewDNSServer(name string) *DNSServer {
	return &DNSServer{
		Name:    name,
		Records: make(map[string]map[DNSRecordType]*DNSRecord),
	}
}

// AddRecord 添加 DNS 记录
func (s *DNSServer) AddRecord(record *DNSRecord) {
	if s.Records[record.Name] == nil {
		s.Records[record.Name] = make(map[DNSRecordType]*DNSRecord)
	}
	s.Records[record.Name][record.Type] = record
}

// Query 查询 DNS 记录
func (s *DNSServer) Query(name string, recordType DNSRecordType) (*DNSRecord, bool) {
	if s.Records[name] == nil {
		return nil, false
	}
	record, exists := s.Records[name][recordType]
	return record, exists
}

// DNSResolver DNS 解析器
type DNSResolver struct {
	LocalCache *DNSCache  // 本地缓存
	LocalDNS   *DNSServer // 本地 DNS 服务器
}

// NewDNSResolver 创建 DNS 解析器
func NewDNSResolver(localDNS *DNSServer) *DNSResolver {
	return &DNSResolver{
		LocalCache: NewDNSCache(),
		LocalDNS:   localDNS,
	}
}

// ResolveRecursive 递归查询
// 对应 408 考点: 客户端向本地 DNS 服务器发起递归查询
func (r *DNSResolver) ResolveRecursive(name string, recordType DNSRecordType) (*DNSRecord, bool) {
	fmt.Printf("\n【递归查询】客户端 → 本地DNS: 查询 %s (%s)\n", name, recordType)

	// 1. 查询本地缓存
	fmt.Printf("  步骤 1: 查询本地缓存...\n")
	if record, found := r.LocalCache.Lookup(name, recordType); found {
		fmt.Printf("  ✓ 命中缓存: %s\n", record)
		return record, true
	}
	fmt.Printf("  ✗ 缓存未命中\n")

	// 2. 向本地 DNS 服务器查询
	fmt.Printf("  步骤 2: 向本地DNS服务器 [%s] 查询...\n", r.LocalDNS.Name)
	if record, found := r.LocalDNS.Query(name, recordType); found {
		fmt.Printf("  ✓ 本地DNS服务器返回: %s\n", record)
		r.LocalCache.Add(record) // 加入缓存
		return record, true
	}

	// 3. 本地 DNS 服务器负责向根、顶级域、权威 DNS 查询 (递归)
	fmt.Printf("  步骤 3: 本地DNS服务器递归查询上级服务器...\n")
	if r.LocalDNS.Parent != nil {
		if record, found := r.LocalDNS.Parent.Query(name, recordType); found {
			fmt.Printf("  ✓ 上级DNS服务器 [%s] 返回: %s\n", r.LocalDNS.Parent.Name, record)
			r.LocalDNS.AddRecord(record) // 本地 DNS 缓存
			r.LocalCache.Add(record)     // 客户端缓存
			return record, true
		}
	}

	fmt.Printf("  ✗ 查询失败: 域名 %s 不存在\n", name)
	return nil, false
}

// ResolveIterative 迭代查询
// 对应 408 考点: DNS 服务器之间的迭代查询
func (r *DNSResolver) ResolveIterative(name string, recordType DNSRecordType) (*DNSRecord, bool) {
	fmt.Printf("\n【迭代查询】客户端主导查询 %s (%s)\n", name, recordType)

	// 1. 查询本地缓存
	fmt.Printf("  步骤 1: 查询本地缓存...\n")
	if record, found := r.LocalCache.Lookup(name, recordType); found {
		fmt.Printf("  ✓ 命中缓存: %s\n", record)
		return record, true
	}
	fmt.Printf("  ✗ 缓存未命中\n")

	// 2. 向本地 DNS 查询
	fmt.Printf("  步骤 2: 向本地DNS [%s] 查询...\n", r.LocalDNS.Name)
	if record, found := r.LocalDNS.Query(name, recordType); found {
		fmt.Printf("  ✓ 返回: %s\n", record)
		r.LocalCache.Add(record)
		return record, true
	}
	fmt.Printf("  ✗ 未找到,返回下一级服务器地址\n")

	// 3. 客户端向上级 DNS 查询 (迭代)
	currentServer := r.LocalDNS.Parent
	step := 3
	for currentServer != nil {
		fmt.Printf("  步骤 %d: 向上级DNS [%s] 查询...\n", step, currentServer.Name)
		if record, found := currentServer.Query(name, recordType); found {
			fmt.Printf("  ✓ 返回: %s\n", record)
			r.LocalCache.Add(record)
			r.LocalDNS.AddRecord(record) // 本地 DNS 学习记录
			return record, true
		}
		fmt.Printf("  ✗ 未找到,继续向上查询\n")
		currentServer = currentServer.Parent
		step++
	}

	fmt.Printf("  ✗ 查询失败: 域名 %s 不存在\n", name)
	return nil, false
}

// DNSExample DNS 协议示例
func DNSExample() {
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("【网络协议 - DNS 域名解析示例】")
	fmt.Println(strings.Repeat("─", 50))

	// 构建 DNS 层次结构
	// 根 DNS → 顶级域 DNS (.com) → 权威 DNS (example.com) → 本地 DNS

	// 1. 创建 DNS 服务器层次
	fmt.Println("\n1️⃣  构建 DNS 服务器层次结构:")

	// 根 DNS 服务器
	rootDNS := NewDNSServer("根DNS服务器")

	// 顶级域 DNS 服务器 (.com)
	comDNS := NewDNSServer("顶级域DNS (.com)")
	comDNS.Parent = rootDNS

	// 权威 DNS 服务器 (example.com)
	exampleDNS := NewDNSServer("权威DNS (example.com)")
	exampleDNS.Parent = comDNS

	// 添加 DNS 记录
	exampleDNS.AddRecord(&DNSRecord{
		Name:  "www.example.com",
		Type:  A,
		Value: "93.184.216.34",
		TTL:   3600,
	})
	exampleDNS.AddRecord(&DNSRecord{
		Name:  "mail.example.com",
		Type:  A,
		Value: "93.184.216.35",
		TTL:   3600,
	})
	exampleDNS.AddRecord(&DNSRecord{
		Name:  "example.com",
		Type:  MX,
		Value: "mail.example.com",
		TTL:   7200,
	})
	exampleDNS.AddRecord(&DNSRecord{
		Name:  "ftp.example.com",
		Type:  CNAME,
		Value: "www.example.com",
		TTL:   3600,
	})

	// 本地 DNS 服务器
	localDNS := NewDNSServer("本地DNS服务器 (ISP)")
	localDNS.Parent = exampleDNS // 简化层次,直接指向权威 DNS

	fmt.Println("  ✓ DNS 层次结构:")
	fmt.Println("      根DNS → 顶级域DNS (.com) → 权威DNS (example.com) → 本地DNS")

	// 2. 创建 DNS 解析器
	resolver := NewDNSResolver(localDNS)

	// 3. 递归查询示例
	fmt.Println("\n2️⃣  DNS 递归查询:")
	fmt.Println(strings.Repeat("═", 50))
	record, found := resolver.ResolveRecursive("www.example.com", A)
	if found {
		fmt.Printf("\n✓ 解析成功: %s → %s\n", record.Name, record.Value)
	}

	// 4. 再次查询同一域名 (命中缓存)
	fmt.Println("\n3️⃣  再次查询 (测试缓存):")
	fmt.Println(strings.Repeat("═", 50))
	record2, found2 := resolver.ResolveRecursive("www.example.com", A)
	if found2 {
		fmt.Printf("\n✓ 解析成功 (缓存): %s → %s\n", record2.Name, record2.Value)
	}

	// 5. 迭代查询示例
	fmt.Println("\n4️⃣  DNS 迭代查询:")
	fmt.Println(strings.Repeat("═", 50))
	record3, found3 := resolver.ResolveIterative("mail.example.com", A)
	if found3 {
		fmt.Printf("\n✓ 解析成功: %s → %s\n", record3.Name, record3.Value)
	}

	// 6. 查询 CNAME 记录
	fmt.Println("\n5️⃣  查询别名 (CNAME 记录):")
	fmt.Println(strings.Repeat("═", 50))
	record4, found4 := resolver.ResolveRecursive("ftp.example.com", CNAME)
	if found4 {
		fmt.Printf("\n✓ 别名解析: %s → %s (CNAME)\n", record4.Name, record4.Value)
		// 进一步解析规范名
		record5, found5 := resolver.ResolveRecursive(record4.Value, A)
		if found5 {
			fmt.Printf("✓ 最终解析: %s → %s (A)\n", record5.Name, record5.Value)
		}
	}

	// 7. 显示 DNS 缓存
	fmt.Println("\n6️⃣  DNS 缓存状态:")
	resolver.LocalCache.PrintCache()

	// 408 考点提示
	fmt.Println("\n📚 408 考点总结:")
	fmt.Println("  ✓ DNS 递归查询: 客户端 → 本地DNS,本地DNS 负责完整解析")
	fmt.Println("  ✓ DNS 迭代查询: DNS 服务器之间,返回下一级服务器地址")
	fmt.Println("  ✓ DNS 记录类型: A (IPv4), AAAA (IPv6), CNAME (别名), MX (邮件), NS (名称服务器)")
	fmt.Println("  ✓ DNS 使用 UDP/53 端口 (查询), TCP/53 端口 (区域传输)")
	fmt.Println("  ✓ DNS 缓存: 减少查询次数,提高解析速度")
	fmt.Println("  ✓ TTL (Time To Live): 记录在缓存中的生存时间")
}
