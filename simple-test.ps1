# 简单的API测试脚本
Write-Host "=== 简单API测试 ==="

# 测试健康检查
Write-Host "\n测试健康检查..."
try {
    $healthResp = Invoke-RestMethod -Uri "http://localhost:9090/healthz" -Method GET
    Write-Host "健康检查成功: 状态码 200" -ForegroundColor Green
} catch {
    Write-Host "健康检查失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 测试电影创建
Write-Host "\n测试电影创建..."
$movieData = @{
    title = "Test Movie 2024"
    releaseDate = "2024-01-01"
    genre = "Action"
}

$jsonData = $movieData | ConvertTo-Json
Write-Host "请求数据: $jsonData"

try {
    $response = Invoke-WebRequest `
        -Uri "http://localhost:9090/movies" `
        -Method POST `
        -Body $jsonData `
        -ContentType "application/json" `
        -Headers @{"Authorization"="Bearer test-token"} `
        -ErrorAction Stop
    
    Write-Host "成功: 状态码 $($response.StatusCode)" -ForegroundColor Green
    Write-Host "响应: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "失败!" -ForegroundColor Red
    if ($_.Exception.Response) {
        Write-Host "HTTP状态码: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
        if ($_.ErrorDetails.Message) {
            Write-Host "错误内容: $($_.ErrorDetails.Message)" -ForegroundColor Red
        }
    } else {
        Write-Host "错误: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host "\n=== 测试结束 ==="