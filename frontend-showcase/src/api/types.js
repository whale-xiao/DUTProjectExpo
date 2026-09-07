/**
 * DTO 契约（与后端 dto/*.go、docs/09 §4.3 对齐，未上 TS 前的类型注释）。
 * 约定：camelCase；数组缺省为 []；分页必有 total/has_more。
 */

/**
 * @typedef {Object} Card
 * @property {number} id
 * @property {string} name
 * @property {string} summary
 * @property {string} coverUrl  空串时前端显示占位
 * @property {string[]} tags
 * @property {boolean} isRecommended
 * @property {number} categoryId
 */

/**
 * @typedef {Object} SuggestItem
 * @property {number} id
 * @property {string} name
 * @property {string} summary
 * @property {string} coverUrl
 * @property {string[]} tags
 */

/**
 * @typedef {Object} Block
 * @property {'heading'|'paragraph'|'image'} type
 * @property {string} [text]
 * @property {string} [url]
 * @property {string} [caption]
 */

/**
 * @typedef {Object} ProjectDetail
 * @property {number} id
 * @property {string} name
 * @property {string} summary
 * @property {string} coverUrl
 * @property {{id:number,name:string}|null} category
 * @property {string[]} tags
 * @property {{name:string,role:string}[]} members
 * @property {number} viewCount
 * @property {string} publishedAt
 * @property {Block[]} blocks
 */

export {}
