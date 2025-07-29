my-microservice/
├── api/                     # API相关定义
│   ├── docs/                # API文档(Swagger等)
│   ├── handlers/            # 请求处理器
│   ├── middleware/          # 自定义中间件
│   └── routes/              # 路由定义
├── cmd/                     # 主程序入口
│   └── main.go              # 主程序文件
├── configs/                 # 配置文件
│   ├── config.yaml          # 主配置文件
│   └── config_dev.yaml      # 开发环境配置
├── internal/                # 内部实现(不对外暴露)
│   ├── models/              # 数据模型/实体
│   ├── repositories/        # 数据访问层
│   ├── services/            # 业务逻辑层
│   └── utils/               # 工具函数
├── pkg/                     # 可复用的公共包
│   ├── database/            # 数据库连接
│   ├── logging/             # 日志处理
│   └── errors/              # 自定义错误
├── scripts/                 # 脚本文件
├── tests/                   # 测试代码
│   ├── integration/         # 集成测试
│   └── unit/                # 单元测试
├── go.mod                   # Go模块定义
├── go.sum                   # Go依赖校验
└── README.md                # 项目说明文档