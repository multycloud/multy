package resources

type MultyContext struct {
	Resources    *Resources
	dependencies map[Resource][]Resource
}

func NewMultyContext(r *Resources) MultyContext {
	return MultyContext{Resources: r, dependencies: map[Resource][]Resource{}}
}

func GetAllResources[T Resource](ctx MultyContext) []T {
	var result []T
	for _, r := range ctx.Resources.ResourceMap {
		if casted, canCast := r.(T); canCast {
			result = append(result, casted)
		}
	}
	return result
}

func GetAllResourcesWithRef[T Resource, T2 Resource](ctx MultyContext, refGetter func(T) T2, ref T2) []T {
	var result []T
	for _, r := range ctx.Resources.ResourceMap {
		if casted, canCast := r.(T); canCast && refGetter(casted).GetResourceId() == ref.GetResourceId() {
			result = append(result, casted)
		}
	}
	return result
}

func GetAllResourcesWithListRef[T Resource, T2 Resource](ctx MultyContext, refGetter func(T) []T2, ref T2) []T {
	var result []T
	for _, r := range ctx.Resources.ResourceMap {
		if casted, canCast := r.(T); canCast {
			for _, tentativeRef := range refGetter(casted) {
				if tentativeRef.GetResourceId() == ref.GetResourceId() {
					result = append(result, casted)
					break
				}
			}

		}
	}
	return result
}
