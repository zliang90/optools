# container stats

# 安装依赖包
pip3 install -r requirements.txt

# 使用

1) 默认排序
> python3 container_stats.py

2) containerd按内存使用量、CPU使用百分比进行排序, 显示前50个容器
> python3 container_stats.py WorkingSetBytes,CPU 50

3) docker按内存使用量、CPU使用百分比进行排序
> python3 container_stats.py MemoryUsed,CPU 50
