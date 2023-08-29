package queue

func Enqueue[V any](queue []V, item V) []V {
	return append(queue, item)
}

func Top[V any](queue []V) V {
	return queue[0]
}

func Pop[V any](queue []V) []V {
	return queue[1:]
}
