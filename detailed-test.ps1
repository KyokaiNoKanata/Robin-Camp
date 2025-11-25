# 详细的API端点测试脚本
Write-Host "=== 详细API端点测试 ===" -ForegroundColor Yellow

# 测试健康检查端点
Write-Host "`n测试健康检查端点..." -ForegroundColor Cyan
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:9090/healthz" -Method GET -ContentType "application/json"
    Write-Host "健康检查成功: $($healthResponse | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "健康检查失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 准备电影创建数据
$movieData = @{
    title = "Test Movie 2024"
    releaseDate = "2024-01-01"
    genre = "Action"
}

Write-Host "`n测试电影创建端点..." -ForegroundColor Cyan
Write-Host "请求数据: $($movieData | ConvertTo-Json)" -ForegroundColor Gray

# 测试电影创建端点（带错误处理）
try {
    # 使用Invoke-WebRequest而不是Invoke-RestMethod以获取更完整的响应
    $response = Invoke-WebRequest -Uri "http://localhost:9090/movies" -Method POST -Body ($movieData | ConvertTo-Json) -ContentType "application/json" -Headers @{"Authorization"="Bearer test-token"} -ErrorAction Stop
    
    Write-Host "电影创建成功! 状态码: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "响应内容: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "电影创建失败!" -ForegroundColor Red
    
    # 获取详细的错误信息
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        $statusDescription = $_.Exception.Response.StatusDescription
        Write-Host "HTTP状态码: $statusCode ($statusDescription)" -ForegroundColor Red
        
        try {
            $errorContent = $_.ErrorDetails.Message
            Write-Host "错误响应内容: $errorContent" -ForegroundColor Red
        } catch {
            Write-Host "无法获取错误响应内容" -ForegroundColor Red
        }
    } else {
        Write-Host "连接错误: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # 显示完整的异常信息
    Write-Host "`n完整异常信息:" -ForegroundColor Magenta
    Write-Host $_.ToString() -ForegroundColor Magenta
    
    if ($_.ScriptStackTrace) {
        Write-Host "`n脚本堆栈跟踪:" -ForegroundColor Magenta
        Write-Host $_.ScriptStackTrace -ForegroundColor Magenta
    }
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Yellow