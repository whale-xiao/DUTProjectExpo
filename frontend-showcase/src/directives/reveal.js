// v-reveal：元素进入视口时渐进浮现，可选错开延迟。见 docs/12 §7。
//
// 为什么用 IntersectionObserver 而不是监听 scroll 事件：
//   scroll 事件每帧触发、跑在主线程上，长列表滚动时容易掉帧；
//   IntersectionObserver 由浏览器在合成线程回调，只在元素进出视口时触发一次。
// 全站共用一个 observer 实例；元素显现后立刻 unobserve（一次性动画），
// 所以滚动再久也不会累积观察开销。
//
// 用法：
//   <div v-reveal>          进入视口时浮现
//   <div v-reveal="120">    延迟 120ms 浮现（网格错开用）
//
// 无障碍与健壮性：
//   - prefers-reduced-motion 下直接标记为已显现，不注册观察、不做过渡
//   - 浏览器不支持 IntersectionObserver 时同样直接显示，绝不让内容卡在隐藏态
//   - .reveal 类由本指令添加，所以 JS 整体失效时内容依然可见

const BASE = 'reveal'
const SHOWN = 'is-revealed'

let observer = null

function prefersReduced() {
  return window.matchMedia?.('(prefers-reduced-motion: reduce)').matches ?? false
}

function getObserver() {
  if (observer) return observer
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue
        const el = entry.target
        const delay = Number(el.dataset.revealDelay || 0)
        if (delay > 0) el.style.transitionDelay = `${delay}ms`
        el.classList.add(SHOWN)
        observer.unobserve(el) // 一次性
      }
    },
    {
      // 露出约 10% 且越过视口下缘 40px 才触发，避免"刚碰到边就弹"
      threshold: 0.1,
      rootMargin: '0px 0px -40px 0px',
    },
  )
  return observer
}

export const vReveal = {
  mounted(el, binding) {
    if (prefersReduced() || typeof IntersectionObserver === 'undefined') {
      el.classList.add(SHOWN)
      return
    }
    el.classList.add(BASE)

    const delay = Number(binding.value) || 0
    if (delay > 0) el.dataset.revealDelay = String(delay)

    getObserver().observe(el)
  },
  unmounted(el) {
    // 只在 observer 已创建时取消；此处不再惰性创建
    observer?.unobserve(el)
  },
}
