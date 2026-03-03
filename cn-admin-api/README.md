省、市、区联动数据查询

### Postman / curl 测试示例
```
# 1) 省列表
curl "http://localhost:8080/api/provinces"

# 2) 某省下的市（例：北京市 110000）
curl "http://localhost:8080/api/cities?province_code=110000"

# 3) 某市下的县/区（例：北京市辖区 110100）
curl "http://localhost:8080/api/counties?city_code=110100"

# 4) 任意节点 children（例：110000 的 children 就是北京市的市级）
curl "http://localhost:8080/api/children/110000"

# 5) 查单节点详情
curl "http://localhost:8080/api/node/110101"

# 6) 模糊搜索（可限定 level）
curl "http://localhost:8080/api/search?q=朝阳&level=3&limit=20"

# 7) 获取路径（省->市->县）
curl "http://localhost:8080/api/path/110101"
```
