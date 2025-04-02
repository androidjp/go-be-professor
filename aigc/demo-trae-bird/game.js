// 游戏配置
const config = {
    canvasWidth: 288,
    canvasHeight: 512,
    gravity: 0.15,
    jumpForce: -4.6,
    pipeGap: 150,
    pipeWidth: 52,
    pipeSpawnInterval: 1500,
    birdSize: 24,
    congratsDuration: 3000 // 祝贺动画持续时间（毫秒）
};

// 获取历史最高分
let highScore = parseInt(localStorage.getItem('flappyBirdHighScore')) || 0;

// 获取画布和上下文
const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');

// 设置画布大小
canvas.width = config.canvasWidth;
canvas.height = config.canvasHeight;

// 游戏状态
let gameStarted = false;
let gameOver = false;
let score = 0;

// 管道生成定时器
let pipeGeneratorTimer = null;

// 小鸟对象
const bird = {
    x: config.canvasWidth / 3,
    y: config.canvasHeight / 2,
    velocity: 0,
    draw() {
        // 保存当前绘图状态
        ctx.save();
        
        // 移动到小鸟的位置
        ctx.translate(this.x, this.y);
        
        // 根据速度稍微旋转小鸟
        const rotation = Math.min(Math.max(this.velocity * 0.1, -0.5), 0.5);
        ctx.rotate(rotation);
        
        // 绘制小鸟身体（黄色圆形）
        ctx.beginPath();
        ctx.fillStyle = '#FFD700';
        ctx.arc(config.birdSize/2, config.birdSize/2, config.birdSize/2, 0, Math.PI * 2);
        ctx.fill();
        
        // 绘制翅膀（深黄色）
        ctx.fillStyle = '#FFA500';
        ctx.beginPath();
        ctx.ellipse(0, config.birdSize/2, config.birdSize/3, config.birdSize/4, Math.PI/4, 0, Math.PI * 2);
        ctx.fill();
        
        // 绘制眼睛（黑色）
        ctx.fillStyle = '#000';
        ctx.beginPath();
        ctx.arc(config.birdSize/1.5, config.birdSize/3, config.birdSize/8, 0, Math.PI * 2);
        ctx.fill();
        
        // 绘制嘴巴（橙色三角形）
        ctx.fillStyle = '#FF6B00';
        ctx.beginPath();
        ctx.moveTo(config.birdSize, config.birdSize/2);
        ctx.lineTo(config.birdSize + config.birdSize/3, config.birdSize/2);
        ctx.lineTo(config.birdSize, config.birdSize/1.5);
        ctx.closePath();
        ctx.fill();
        
        // 恢复绘图状态
        ctx.restore();
    },
    update() {
        if (gameStarted) {
            this.velocity += config.gravity;
            this.y += this.velocity;
        }
    },
    jump() {
        this.velocity = config.jumpForce;
    },
    reset() {
        this.y = config.canvasHeight / 2;
        this.velocity = 0;
    }
};

// 管道数组
let pipes = [];

// 管道类
class Pipe {
    constructor() {
        this.x = config.canvasWidth;
        this.gapY = Math.random() * (config.canvasHeight - config.pipeGap - 100) + 50;
        this.passed = false;
    }

    draw() {
        // 创建渐变色
        const gradient = ctx.createLinearGradient(this.x, 0, this.x + config.pipeWidth, 0);
        gradient.addColorStop(0, '#8B4513');
        gradient.addColorStop(0.5, '#CD853F');
        gradient.addColorStop(1, '#8B4513');

        // 上管道
        ctx.fillStyle = gradient;
        ctx.fillRect(this.x, 0, config.pipeWidth, this.gapY);

        // 上管道装饰边缘
        ctx.fillStyle = '#A0522D';
        ctx.fillRect(this.x - 3, this.gapY - 20, config.pipeWidth + 6, 20);
        
        // 上管道纹理
        ctx.strokeStyle = '#6B4423';
        ctx.lineWidth = 1;
        for(let i = 20; i < this.gapY - 20; i += 20) {
            ctx.beginPath();
            ctx.moveTo(this.x, i);
            ctx.lineTo(this.x + config.pipeWidth, i);
            ctx.stroke();
        }

        // 下管道
        ctx.fillStyle = gradient;
        ctx.fillRect(
            this.x,
            this.gapY + config.pipeGap,
            config.pipeWidth,
            config.canvasHeight - (this.gapY + config.pipeGap)
        );

        // 下管道装饰边缘
        ctx.fillStyle = '#A0522D';
        ctx.fillRect(this.x - 3, this.gapY + config.pipeGap, config.pipeWidth + 6, 20);

        // 下管道纹理
        for(let i = this.gapY + config.pipeGap + 40; i < config.canvasHeight - 20; i += 20) {
            ctx.beginPath();
            ctx.moveTo(this.x, i);
            ctx.lineTo(this.x + config.pipeWidth, i);
            ctx.stroke();
        }

        // 添加光影效果
        ctx.fillStyle = 'rgba(255, 255, 255, 0.1)';
        ctx.fillRect(this.x, 0, 5, this.gapY);
        ctx.fillRect(this.x, this.gapY + config.pipeGap, 5, config.canvasHeight - (this.gapY + config.pipeGap));
    }

    update() {
        this.x -= 2;
        // 检查是否通过管道
        if (!this.passed && this.x + config.pipeWidth < bird.x) {
            score++;
            this.passed = true;
        }
    }

    checkCollision() {
        const birdRight = bird.x + config.birdSize;
        const birdBottom = bird.y + config.birdSize;
        
        // 检查是否碰到上管道
        if (bird.x < this.x + config.pipeWidth &&
            birdRight > this.x &&
            bird.y < this.gapY) {
            return true;
        }
        
        // 检查是否碰到下管道
        if (bird.x < this.x + config.pipeWidth &&
            birdRight > this.x &&
            birdBottom > this.gapY + config.pipeGap) {
            return true;
        }
        
        return false;
    }
}

// 生成管道
function spawnPipe() {
    if (gameStarted && !gameOver) {
        pipes.push(new Pipe());
    }
}

// 绘制背景
function drawBackground() {
    // 天空 - 根据分数切换昼夜
    let skyGradient;
    if (score >= 5) {
        // 夜晚天空
        skyGradient = ctx.createLinearGradient(0, 0, 0, config.canvasHeight);
        skyGradient.addColorStop(0, '#0B1026');
        skyGradient.addColorStop(1, '#2B4570');
    } else {
        // 白天天空
        skyGradient = ctx.createLinearGradient(0, 0, 0, config.canvasHeight);
        skyGradient.addColorStop(0, '#4ec0ca');
        skyGradient.addColorStop(1, '#87CEEB');
    }
    ctx.fillStyle = skyGradient;
    ctx.fillRect(0, 0, config.canvasWidth, config.canvasHeight);

    // 如果是夜晚，绘制星星和月亮
    if (score >= 5) {
        // 绘制星星
        ctx.fillStyle = '#FFF';
        const stars = [
            {x: 50, y: 50},
            {x: 100, y: 80},
            {x: 150, y: 30},
            {x: 200, y: 70},
            {x: 250, y: 40},
            {x: 30, y: 100},
            {x: 80, y: 150},
            {x: 130, y: 120},
            {x: 180, y: 90},
            {x: 230, y: 130}
        ];
        stars.forEach(star => {
            const size = 2;
            ctx.beginPath();
            ctx.arc(star.x, star.y, size, 0, Math.PI * 2);
            ctx.fill();
        });

        // 绘制月亮
        ctx.beginPath();
        ctx.fillStyle = '#FFE5B4';
        ctx.arc(config.canvasWidth - 60, 60, 30, 0, Math.PI * 2);
        ctx.fill();
        
        // 月亮的阴影
        ctx.beginPath();
        ctx.fillStyle = 'rgba(255, 229, 180, 0.3)';
        ctx.arc(config.canvasWidth - 55, 55, 35, 0, Math.PI * 2);
        ctx.fill();
    }

    // 绘制云朵
    function drawCloud(x, y, scale) {
        ctx.save();
        ctx.translate(x, y);
        ctx.scale(scale, scale);

        // 创建云朵渐变
        const cloudGradient = ctx.createRadialGradient(0, 0, 0, 0, 0, 40);
        cloudGradient.addColorStop(0, 'rgba(255, 255, 255, 0.9)');
        cloudGradient.addColorStop(1, 'rgba(255, 255, 255, 0.2)');

        // 绘制主体
        ctx.fillStyle = cloudGradient;
        ctx.beginPath();
        ctx.arc(0, 0, 20, 0, Math.PI * 2);
        ctx.arc(15, -10, 15, 0, Math.PI * 2);
        ctx.arc(-15, -8, 18, 0, Math.PI * 2);
        ctx.arc(-8, 8, 16, 0, Math.PI * 2);
        ctx.arc(10, 5, 17, 0, Math.PI * 2);
        ctx.fill();

        ctx.restore();
    }

    // 绘制多个云朵
    const time = Date.now() * 0.001;
    drawCloud(50 + Math.sin(time) * 10, 80, 1);
    drawCloud(150 + Math.cos(time * 0.8) * 15, 150, 1.2);
    drawCloud(250 + Math.sin(time * 1.2) * 12, 100, 0.8);
    
    // 草地 - 使用渐变色
    const grassGradient = ctx.createLinearGradient(0, config.canvasHeight - 50, 0, config.canvasHeight);
    grassGradient.addColorStop(0, '#5EBA7D');
    grassGradient.addColorStop(1, '#3E8A5D');
    ctx.fillStyle = grassGradient;
    ctx.fillRect(0, config.canvasHeight - 50, config.canvasWidth, 50);
    
    // 添加草地纹理
    ctx.strokeStyle = '#4EA86C';
    ctx.lineWidth = 1;
    for(let i = 0; i < config.canvasWidth; i += 15) {
        const height = Math.random() * 8 + 5;
        ctx.beginPath();
        ctx.moveTo(i, config.canvasHeight - 50);
        ctx.lineTo(i + 5, config.canvasHeight - 50 + height);
        ctx.lineTo(i + 10, config.canvasHeight - 50);
        ctx.stroke();
    }
    
    // 城市轮廓
    const buildingColors = score >= 100 ? ['#1a1a1a', '#2B2B2B', '#3a3a3a'] : ['#DDDDDD', '#CCCCCC', '#BBBBBB'];
    const buildingCount = Math.floor(config.canvasWidth / 40);
    
    for (let i = 0; i < buildingCount; i++) {
        const x = i * 40;
        const height = 30 + Math.sin(i * 0.8) * 20;
        const buildingWidth = 35;
        
        // 绘制建筑主体
        ctx.fillStyle = buildingColors[i % buildingColors.length];
        ctx.fillRect(x, config.canvasHeight - 100 - height, buildingWidth, height + 50);
        
        // 添加建筑顶部装饰
        if (i % 3 === 0) {
            ctx.beginPath();
            ctx.moveTo(x, config.canvasHeight - 100 - height);
            ctx.lineTo(x + buildingWidth/2, config.canvasHeight - 110 - height);
            ctx.lineTo(x + buildingWidth, config.canvasHeight - 100 - height);
            ctx.fillStyle = score >= 100 ? '#4a4a4a' : '#AAAAAA';
            ctx.fill();
        }
        
        // 添加窗户
        const windowRows = Math.floor(height / 15);
        const windowCols = 3;
        const windowSize = 8;
        const windowSpacing = (buildingWidth - windowCols * windowSize) / (windowCols + 1);
        
        for (let row = 0; row < windowRows; row++) {
            for (let col = 0; col < windowCols; col++) {
                const windowX = x + windowSpacing + col * (windowSize + windowSpacing);
                const windowY = config.canvasHeight - 90 - height + row * 15;
                
                if (score >= 5) {
                    // 夜晚发光窗户效果
                    const glowColor = Math.random() > 0.3 ? 'rgba(255, 255, 150, 0.8)' : 'rgba(255, 255, 150, 0.2)';
                    ctx.fillStyle = glowColor;
                    
                    // 添加光晕效果
                    const glow = ctx.createRadialGradient(
                        windowX + windowSize/2, windowY + windowSize/2, 0,
                        windowX + windowSize/2, windowY + windowSize/2, windowSize
                    );
                    glow.addColorStop(0, glowColor);
                    glow.addColorStop(1, 'rgba(255, 255, 150, 0)');
                    ctx.fillStyle = glow;
                    ctx.fillRect(windowX - 2, windowY - 2, windowSize + 4, windowSize + 4);
                    
                    // 绘制窗户
                    ctx.fillStyle = glowColor;
                    ctx.fillRect(windowX, windowY, windowSize, windowSize);
                } else {
                    // 白天窗户效果
                    ctx.fillStyle = 'rgba(200, 200, 200, 0.9)';
                    ctx.fillRect(windowX, windowY, windowSize, windowSize);
                    
                    // 添加窗户反光效果
                    ctx.fillStyle = 'rgba(255, 255, 255, 0.4)';
                    ctx.fillRect(windowX, windowY, windowSize/2, windowSize/2);
                }
            }
        }
    }
}

// 绘制分数
function drawScore() {
    ctx.fillStyle = '#000';
    ctx.font = '24px Arial';
    ctx.fillText(`分数: ${score}`, 10, 30);
    ctx.fillText(`最高分: ${highScore}`, 10, 60);
}

// 绘制奖牌和祝贺文字
function drawCongrats() {
    // 绘制半透明背景
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
    ctx.fillRect(0, 0, config.canvasWidth, config.canvasHeight);
    
    // 绘制奖牌
    ctx.beginPath();
    ctx.arc(config.canvasWidth / 2, config.canvasHeight / 2 - 40, 30, 0, Math.PI * 2);
    const medalGradient = ctx.createRadialGradient(
        config.canvasWidth / 2, config.canvasHeight / 2 - 40, 0,
        config.canvasWidth / 2, config.canvasHeight / 2 - 40, 30
    );
    medalGradient.addColorStop(0, '#FFD700');
    medalGradient.addColorStop(1, '#FFA500');
    ctx.fillStyle = medalGradient;
    ctx.fill();
    ctx.strokeStyle = '#B8860B';
    ctx.lineWidth = 3;
    ctx.stroke();
    
    // 绘制星星
    ctx.fillStyle = '#FFFFFF';
    ctx.font = '30px Arial';
    ctx.textAlign = 'center';
    ctx.fillText('★', config.canvasWidth / 2, config.canvasHeight / 2 - 30);
    
    // 绘制祝贺文字
    ctx.font = 'bold 32px Arial';
    const congratsGradient = ctx.createLinearGradient(
        config.canvasWidth / 2 - 80,
        config.canvasHeight / 2 + 20,
        config.canvasWidth / 2 + 80,
        config.canvasHeight / 2 + 20
    );
    congratsGradient.addColorStop(0, '#FF0000');
    congratsGradient.addColorStop(0.5, '#FFFF00');
    congratsGradient.addColorStop(1, '#FF0000');
    ctx.fillStyle = congratsGradient;
    ctx.fillText('你真棒！', config.canvasWidth / 2, config.canvasHeight / 2 + 20);
    
    // 绘制新纪录文字
    ctx.font = '24px Arial';
    ctx.fillStyle = '#FFFFFF';
    ctx.fillText('新纪录：' + score, config.canvasWidth / 2, config.canvasHeight / 2 + 60);
}

// 绘制游戏开始/结束提示
function drawMessage() {
    ctx.fillStyle = '#000';
    ctx.font = '20px Arial';
    if (!gameStarted) {
        ctx.fillText('点击开始游戏', config.canvasWidth / 2 - 60, config.canvasHeight / 2);
    } else if (gameOver) {
        ctx.fillText('游戏结束 - 点击重新开始', config.canvasWidth / 2 - 100, config.canvasHeight / 2);
    }
}

// 游戏主循环
function gameLoop() {
    // 清空画布
    ctx.clearRect(0, 0, config.canvasWidth, config.canvasHeight);
    
    // 绘制背景
    drawBackground();

    // 更新和绘制管道
    pipes = pipes.filter(pipe => pipe.x + config.pipeWidth > 0);
    pipes.forEach(pipe => {
        pipe.update();
        pipe.draw();
        if (pipe.checkCollision()) {
            gameOver = true;
        }
    });

    // 更新和绘制小鸟
    bird.update();
    bird.draw();

    // 检查小鸟是否撞到地面或飞出顶部
    if (bird.y + config.birdSize > config.canvasHeight - 50 || bird.y < 0) {
        gameOver = true;
    }

    // 绘制游戏消息
    drawMessage();
    
    // 绘制分数（最后绘制以确保显示在最上层）
    ctx.save();
    ctx.globalCompositeOperation = 'source-over';
    drawScore();
    ctx.restore();

    // 检查游戏结束
    if (gameOver) {
        // 检查是否打破记录
        if (score > highScore) {
            highScore = score;
            localStorage.setItem('flappyBirdHighScore', highScore);
            drawCongrats();
            setTimeout(() => {
                // 重置游戏
                gameStarted = false;
                gameOver = false;
                score = 0;
                pipes = [];
                bird.reset();
                // 清除画布并重新绘制背景和分数
                ctx.clearRect(0, 0, config.canvasWidth, config.canvasHeight);
                drawBackground();
                bird.draw();
                drawScore();
                drawMessage();
                setTimeout(() => {
                    gameStarted = true;
                    gameLoop();
                }, 100);
            }, config.congratsDuration);
        }
    }

    // 继续游戏循环
    if (!gameOver) {
        requestAnimationFrame(gameLoop);
    }
}

// 事件监听
canvas.addEventListener('click', () => {
    if (!gameStarted) {
        gameStarted = true;
        if (pipeGeneratorTimer) clearInterval(pipeGeneratorTimer);
        pipeGeneratorTimer = setInterval(spawnPipe, config.pipeSpawnInterval);
        gameLoop();
    } else if (gameOver) {
        // 重置游戏
        gameStarted = false;
        gameOver = false;
        score = 0;
        pipes = [];
        bird.reset();
        setTimeout(() => {
            gameStarted = true;
            gameLoop();
        }, 100);
    } else {
        bird.jump();
    }
});

// 添加空格键事件监听
document.addEventListener('keydown', (event) => {
    if (event.code === 'Space') {
        event.preventDefault(); // 防止页面滚动
        if (!gameStarted) {
            gameStarted = true;
            if (pipeGeneratorTimer) clearInterval(pipeGeneratorTimer);
            pipeGeneratorTimer = setInterval(spawnPipe, config.pipeSpawnInterval);
            gameLoop();
        } else if (gameOver) {
            // 重置游戏
            gameStarted = false;
            gameOver = false;
            score = 0;
            pipes = [];
            bird.reset();
            setTimeout(() => {
                gameStarted = true;
                gameLoop();
            }, 100);
        } else {
            bird.jump();
        }
    }
});

// 初始化游戏
drawBackground();
bird.draw();
drawMessage();