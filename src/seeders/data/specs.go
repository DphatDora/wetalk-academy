package data

func lesson(title, objective, language, code string, quizzes int) LessonSpec {
	return LessonSpec{Title: title, Objective: objective, Language: language, Code: code, QuizCount: quizzes, Material: materialFor(title, objective)}
}

func topicSpecs() []TopicSpec {
	return []TopicSpec{
		{
			Key: "go-programming", Title: "Go Programming Căn Bản", Language: "go",
			Description: "Nắm nền tảng Go và các công cụ cần thiết để viết chương trình đồng thời, dễ kiểm thử và dễ bảo trì.",
			Lessons: []LessonSpec{
				lesson("Môi trường Go và module", "Khởi tạo module, quản lý dependency và chạy chương trình Go", "bash", "go mod init example/app\ngo run ./cmd/app", 1),
				lesson("Kiểu dữ liệu và collection", "Sử dụng type, slice và map an toàn", "go", "scores := map[string]int{\"An\": 9}\nscores[\"Binh\"] = 8", 1),
				lesson("Control flow và function", "Tổ chức luồng điều khiển và function có lỗi trả về", "go", "func divide(a, b int) (int, error) {\n if b == 0 { return 0, errors.New(\"zero\") }\n return a / b, nil\n}", 1),
				lesson("Struct và interface", "Thiết kế abstraction nhỏ bằng struct và interface", "go", "type Store interface { Save(context.Context, Item) error }", 1),
				lesson("Error handling và testing", "Bao bọc lỗi và viết table-driven test", "go", "if err != nil { return fmt.Errorf(\"save item: %w\", err) }", 1),
				lesson("Goroutine, channel và context", "Điều phối concurrency và hủy công việc đúng cách", "go", "select {\ncase result := <-results:\n return result\ncase <-ctx.Done():\n return ctx.Err()\n}", 2),
			},
		},
		{
			Key: "backend-go", Title: "Backend API Với Go", Language: "go",
			Description: "Xây dựng REST API bằng Go, Gin, MongoDB cùng các nguyên tắc validation, authentication và deployment.",
			Lessons: []LessonSpec{
				lesson("HTTP và Gin router", "Thiết kế route và handler HTTP rõ ràng", "go", "r := gin.New()\nr.GET(\"/health\", healthHandler)", 1),
				lesson("Thiết kế REST API", "Thiết kế resource, status code và response nhất quán", "json", "{\"success\":true,\"data\":{\"id\":\"123\"}}", 1),
				lesson("Validation và error handling", "Xác thực request và chuẩn hóa lỗi API", "go", "if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, errResponse(err)); return }", 1),
				lesson("Middleware và JWT", "Xác thực JWT và truyền identity qua request context", "go", "authorized.Use(AuthMiddleware(conf))", 1),
				lesson("MongoDB repository", "Tách persistence bằng repository và truy vấn MongoDB có index", "go", "collection.FindOne(ctx, bson.M{\"slug\": slug}).Decode(&topic)", 1),
				lesson("Testing và deployment", "Kiểm thử handler và build dịch vụ có thể triển khai", "bash", "go test ./...\nCGO_ENABLED=0 go build ./cmd/server", 1),
			},
		},
		{
			Key: "modern-js-ts", Title: "JavaScript Và TypeScript Hiện Đại", Language: "typescript",
			Description: "Làm chủ runtime JavaScript, bất đồng bộ, type system, tooling, kiểm thử và bảo mật frontend.",
			Lessons: []LessonSpec{
				lesson("JavaScript runtime", "Hiểu event loop, scope và cơ chế thực thi JavaScript", "javascript", "queueMicrotask(() => console.log('microtask'));\nconsole.log('sync');", 1),
				lesson("Async programming", "Phối hợp Promise, async/await và xử lý lỗi bất đồng bộ", "typescript", "const [user, orders] = await Promise.all([getUser(), getOrders()]);", 1),
				lesson("TypeScript type system", "Mô hình hóa dữ liệu bằng union và generic", "typescript", "type Result<T> = { ok: true; value: T } | { ok: false; error: Error };", 1),
				lesson("Module và tooling", "Tổ chức module và cấu hình công cụ build", "json", "{\"compilerOptions\":{\"strict\":true,\"module\":\"ESNext\"}}", 1),
				lesson("Testing JavaScript", "Viết unit test tập trung vào hành vi", "typescript", "expect(calculateTotal(items)).toBe(120);", 1),
				lesson("Performance và security", "Đo hiệu năng và giảm rủi ro XSS, supply-chain", "typescript", "element.textContent = untrustedInput;", 1),
			},
		},
		{
			Key: "react-engineering", Title: "React Frontend Engineering", Language: "typescript",
			Description: "Xây dựng ứng dụng React có cấu trúc tốt, truy cập được, hiệu năng ổn định và dễ kiểm thử.",
			Lessons: []LessonSpec{
				lesson("Component và state", "Phân tách component và quản lý state tối thiểu", "tsx", "function Counter(){ const [count,setCount]=useState(0); return <button onClick={()=>setCount(count+1)}>{count}</button> }", 1),
				lesson("Hooks và data fetching", "Quản lý lifecycle và trạng thái tải dữ liệu", "tsx", "const { data, error, isLoading } = useQuery({ queryKey:['topics'], queryFn:getTopics });", 1),
				lesson("Form validation", "Xây dựng form có validation và phản hồi lỗi rõ ràng", "tsx", "if (!email.includes('@')) setError('Email không hợp lệ');", 1),
				lesson("Routing và authentication", "Bảo vệ route và xử lý phiên đăng nhập", "tsx", "return user ? <Outlet /> : <Navigate to=\"/login\" replace />;", 1),
				lesson("Performance và accessibility", "Tối ưu render đồng thời bảo đảm accessibility", "tsx", "<button aria-label=\"Đóng hộp thoại\" onClick={onClose}>×</button>", 1),
				lesson("Testing và deployment React", "Kiểm thử luồng người dùng và tạo production build", "bash", "npm test\nnpm run build", 0),
			},
		},
		{
			Key: "python-automation-data", Title: "Python Cho Automation Và Data", Language: "python",
			Description: "Sử dụng Python để tự động hóa, xử lý dữ liệu, xây CLI và đóng gói chương trình tin cậy.",
			Lessons: []LessonSpec{
				lesson("Python environment", "Quản lý virtual environment và dependency", "bash", "python -m venv .venv\npython -m pip install -r requirements.txt", 1),
				lesson("File và data structure", "Đọc file an toàn và xử lý collection Python", "python", "with open('input.txt', encoding='utf-8') as f:\n    lines = [line.strip() for line in f]", 1),
				lesson("Gọi API và scraping", "Gọi HTTP API có timeout và xử lý lỗi", "python", "response = requests.get(url, timeout=10)\nresponse.raise_for_status()", 1),
				lesson("CLI automation", "Xây CLI tự động hóa có tham số rõ ràng", "python", "parser.add_argument('--dry-run', action='store_true')", 1),
				lesson("Pandas căn bản", "Làm sạch và tổng hợp dữ liệu dạng bảng", "python", "summary = df.dropna().groupby('category')['amount'].sum()", 1),
				lesson("Testing và packaging", "Kiểm thử và đóng gói project Python", "bash", "pytest -q\npython -m build", 0),
			},
		},
		{
			Key: "devops-cicd", Title: "DevOps Và CI/CD Thực Chiến", Language: "yaml",
			Description: "Tự động hóa build, test, đóng gói, triển khai và quan sát dịch vụ.",
			Lessons: []LessonSpec{
				lesson("Linux và Bash", "Viết script shell có kiểm soát lỗi", "bash", "set -euo pipefail\nprintf 'deploying %s\\n' \"$VERSION\"", 1),
				lesson("Docker image", "Tạo image nhỏ và chạy process không đặc quyền", "dockerfile", "FROM alpine:3.20\nUSER 10001\nCOPY app /app\nENTRYPOINT [\"/app\"]", 1),
				lesson("Docker Compose", "Điều phối dịch vụ local bằng Compose", "yaml", "services:\n  api:\n    build: .\n    depends_on: [mongo]", 1),
				lesson("GitHub Actions", "Thiết kế pipeline CI chạy test và build", "yaml", "steps:\n  - uses: actions/checkout@v4\n  - run: go test ./...", 1),
				lesson("Deployment và observability", "Triển khai có health check, log và metric", "yaml", "readinessProbe:\n  httpGet:\n    path: /api/v1/health\n    port: 8046", 1),
				lesson("Kubernetes căn bản", "Khai báo workload và service Kubernetes", "yaml", "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: academy", 2),
			},
		},
		{
			Key: "database-engineering", Title: "Database Engineering Với MongoDB Và SQL", Language: "sql",
			Description: "Thiết kế dữ liệu, index, transaction, migration và tối ưu truy vấn cho hệ thống thực tế.",
			Lessons: []LessonSpec{
				lesson("Data modeling", "Chọn mô hình nhúng hoặc tham chiếu theo access pattern", "json", "{\"lesson_id\":\"...\",\"sections\":[{\"type\":\"text\"}]}", 1),
				lesson("CRUD và index MongoDB", "Tạo index dựa trên truy vấn thực tế", "javascript", "db.lessons.createIndex({ topic_id: 1, order_index: 1 }, { unique: true })", 1),
				lesson("Relational design", "Chuẩn hóa bảng và xác định khóa ngoại", "sql", "CREATE TABLE lessons (id UUID PRIMARY KEY, topic_id UUID NOT NULL REFERENCES topics(id));", 1),
				lesson("Transaction và consistency", "Lựa chọn transaction và consistency phù hợp", "sql", "BEGIN;\nUPDATE accounts SET balance = balance - 100 WHERE id = 1;\nCOMMIT;", 1),
				lesson("Migration và backup", "Thực hiện migration và backup có thể phục hồi", "bash", "mongodump --uri \"$MONGO_URI\" --out backup", 1),
				lesson("Query optimization", "Đọc execution plan và tối ưu truy vấn", "sql", "EXPLAIN ANALYZE SELECT * FROM lessons WHERE topic_id = $1 ORDER BY order_index;", 1),
			},
		},
		{
			Key: "system-design", Title: "System Design Căn Bản", Language: "text",
			Description: "Phân tích yêu cầu và thiết kế hệ thống có khả năng mở rộng, tin cậy và quan sát được.",
			Lessons: []LessonSpec{
				lesson("Requirement và capacity", "Chuyển yêu cầu sản phẩm thành SLO và ước lượng tải", "text", "DAU=100000, requests/user/day=20, peak_factor=5", 1),
				lesson("API và data model", "Thiết kế API cùng data model từ access pattern", "http", "GET /api/v1/topics/{slug}/lessons?page=1&limit=20", 1),
				lesson("Cache và CDN", "Chọn cache key, TTL và chiến lược invalidation", "text", "cache-key: topic:{slug}:lessons:v1; ttl: 300s", 1),
				lesson("Queue và xử lý bất đồng bộ", "Thiết kế consumer idempotent và retry có kiểm soát", "json", "{\"event_id\":\"uuid\",\"type\":\"lesson.published\",\"version\":1}", 1),
				lesson("Scalability và reliability", "Thiết kế scale-out, redundancy và graceful degradation", "text", "SLO 99.9%; timeout 2s; retry budget 1", 1),
				lesson("Thiết kế academy platform", "Kết hợp API, database, cache, queue và observability", "text", "client -> api gateway -> academy service -> mongo/cache/event bus", 2),
			},
		},
		{
			Key: "cloud-native-microservices", Title: "Cloud Native Và Microservices", Language: "yaml",
			Description: "Thiết kế ranh giới dịch vụ, contract, messaging, resilience, observability và deployment.",
			Lessons: []LessonSpec{
				lesson("Service boundaries", "Xác định bounded context và quyền sở hữu dữ liệu", "text", "Academy owns topics, lessons, contents and quizzes.", 1),
				lesson("Contract và configuration", "Quản lý API contract và cấu hình theo môi trường", "yaml", "database:\n  name: wetalk_academy\nserver:\n  url: https://api.example.com", 1),
				lesson("Messaging", "Thiết kế event có version và consumer idempotent", "json", "{\"id\":\"event-1\",\"type\":\"quiz.submitted\",\"version\":1}", 1),
				lesson("Resilience patterns", "Áp dụng timeout, retry, circuit breaker đúng chỗ", "text", "timeout=2s; max_retries=2; exponential_backoff=true", 1),
				lesson("Observability", "Liên kết log, metric và trace để điều tra sự cố", "json", "{\"trace_id\":\"abc\",\"level\":\"error\",\"message\":\"mongo timeout\"}", 1),
				lesson("Security và deployment", "Triển khai least privilege và quản lý secret", "yaml", "securityContext:\n  runAsNonRoot: true\n  readOnlyRootFilesystem: true", 1),
			},
		},
		{
			Key: "web-security", Title: "Web Security Cho Developer", Language: "text",
			Description: "Nhận diện và giảm thiểu rủi ro phổ biến trong ứng dụng web và API.",
			Lessons: []LessonSpec{
				lesson("Threat modeling và OWASP", "Xác định asset, trust boundary và threat trước khi code", "text", "asset: access token; boundary: browser -> API; threat: token theft", 1),
				lesson("Session và JWT", "Xác thực token, expiry và quyền truy cập", "text", "validate signature, issuer, audience, expiry and subject", 1),
				lesson("Injection và input validation", "Ngăn injection bằng validation và parameterized query", "sql", "SELECT * FROM users WHERE email = $1;", 1),
				lesson("Secret và cryptography", "Quản lý secret và dùng primitive mật mã đúng cách", "bash", "export JWT_SECRET=\"$(openssl rand -base64 32)\"", 1),
				lesson("API security và CORS", "Giới hạn origin, rate và dữ liệu trả về", "http", "Access-Control-Allow-Origin: https://app.example.com", 1),
				lesson("Monitoring và incident response", "Phát hiện, cô lập và phục hồi sau sự cố", "text", "detect -> contain -> eradicate -> recover -> review", 0),
			},
		},
	}
}
