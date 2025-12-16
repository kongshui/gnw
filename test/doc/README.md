
1. 裸库迁移
   在本地可以通过下列命令迁移 全部分支和 TAG 到本仓库，点击创建 访问令牌
    mkdir empty && cd empty
    git clone --bare https://your-git.com/group/name.git .
    git lfs fetch origin --all
    git push --mirror https://cnb.cool/weifenggame/kun_swallows_world.git

2. 分支迁移
   在本地可以通过下列命令迁移 当前分支 到本仓库，点击创建 访问令牌
   git remote -v
   git remote add cnb https://cnb.cool/weifenggame/kun_swallows_world.git
   git remote -v
   git pull origin main --allow-unrelated-histories
   <!-- git push -u cnb HEAD:main -->

3. 空仓初始化
   在本地可以通过下列命令创建 全新仓库, 点击创建 访问令牌
    git init .
    git remote add origin https://cnb.cool/weifenggame/kun_swallows_world.git
    git config --local user.name cnb.aaaaaa    // 自己的名称
    git config --local user.email "your.email@example.com"  //自己的邮箱
    git config credential.helper store
   