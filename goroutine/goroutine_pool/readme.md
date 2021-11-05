# 协程池

协程池内的一定数量的协程。

任务队列，即jobCh，存在协程池不能立即处理任务的情况，所以需要队列把任务先暂存。

结果队列，即retCh，同上，协程池处理任务的结果，也存在不能被下游立刻提取的情况，要暂时保存。

协程池最简要（核心）的逻辑是所有协程从任务读取任务，处理后把结果存放到结果队列。