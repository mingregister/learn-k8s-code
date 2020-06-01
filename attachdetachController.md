https://mp.weixin.qq.com/s/CElaEP5YCR3QF4ppOOicAw

mingregister-attachdetachController

kubernetes-1.18.2\pkg\controller\volume\attachdetach\cache\actual_state_of_world.go
kubernetes-1.18.2\pkg\controller\volume\attachdetach\cache\desired_state_of_world.go

对于 ad controller 来说，理解了其内部的数据结构，再去理解逻辑就事半功倍。ad controller 在内存中维护 2 个数据结构：
* actualStateOfWorld —— 表征实际状态（后面简称 asw）
* desiredStateOfWorld —— 表征期望状态（后面简称 dsw）