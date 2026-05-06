function _createForOfIteratorHelper(r, e) { var t = "undefined" != typeof Symbol && r[Symbol.iterator] || r["@@iterator"]; if (!t) { if (Array.isArray(r) || (t = _unsupportedIterableToArray(r)) || e && r && "number" == typeof r.length) { t && (r = t); var _n = 0, F = function F() {}; return { s: F, n: function n() { return _n >= r.length ? { done: !0 } : { done: !1, value: r[_n++] }; }, e: function e(r) { throw r; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var o, a = !0, u = !1; return { s: function s() { t = t.call(r); }, n: function n() { var r = t.next(); return a = r.done, r; }, e: function e(r) { u = !0, o = r; }, f: function f() { try { a || null == t.return || t.return(); } finally { if (u) throw o; } } }; }
function _unsupportedIterableToArray(r, a) { if (r) { if ("string" == typeof r) return _arrayLikeToArray(r, a); var t = {}.toString.call(r).slice(8, -1); return "Object" === t && r.constructor && (t = r.constructor.name), "Map" === t || "Set" === t ? Array.from(r) : "Arguments" === t || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(t) ? _arrayLikeToArray(r, a) : void 0; } }
function _arrayLikeToArray(r, a) { (null == a || a > r.length) && (a = r.length); for (var e = 0, n = Array(a); e < a; e++) n[e] = r[e]; return n; }
function _typeof(o) { "@babel/helpers - typeof"; return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (o) { return typeof o; } : function (o) { return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o; }, _typeof(o); }
var arr = [];
var each = arr.forEach;
var slice = arr.slice;
var UNSAFE_KEYS = ['__proto__', 'constructor', 'prototype'];
export function defaults(obj) {
  each.call(slice.call(arguments, 1), function (source) {
    if (source) {
      for (var _i = 0, _Object$keys = Object.keys(source); _i < _Object$keys.length; _i++) {
        var prop = _Object$keys[_i];
        if (UNSAFE_KEYS.indexOf(prop) > -1) continue;
        if (obj[prop] === undefined) obj[prop] = source[prop];
      }
    }
  });
  return obj;
}
function isSafeUrlSegmentBase(v) {
  if (typeof v !== 'string') return false;
  if (v.length === 0 || v.length > 128) return false;
  if (UNSAFE_KEYS.indexOf(v) > -1) return false;
  if (v.indexOf('..') > -1) return false;
  if (v.indexOf('\\') > -1) return false;
  if (/[?#%\s@]/.test(v)) return false;
  if (/[\x00-\x1F\x7F]/.test(v)) return false;
  return true;
}
export function isSafeLangUrlSegment(v) {
  if (!isSafeUrlSegmentBase(v)) return false;
  if (v.indexOf('/') > -1) return false;
  return true;
}
export function isSafeNsUrlSegment(v) {
  return isSafeUrlSegmentBase(v);
}
export var isSafeUrlSegment = isSafeLangUrlSegment;
var SAFETY_CHECK_BY_KEY = {
  lng: isSafeLangUrlSegment,
  ns: isSafeNsUrlSegment
};
export function sanitizeLogValue(v) {
  if (typeof v !== 'string') return v;
  return v.replace(/[\r\n\x00-\x1F\x7F]/g, ' ');
}
export function redactUrlCredentials(u) {
  if (typeof u !== 'string' || u.length === 0) return u;
  try {
    var parsed = new URL(u);
    if (parsed.username || parsed.password) {
      parsed.username = '';
      parsed.password = '';
      return parsed.toString();
    }
    return u;
  } catch (e) {
    return u.replace(/(\/\/)[^/@\s]+@/g, '$1');
  }
}
export function hasXMLHttpRequest() {
  return typeof XMLHttpRequest === 'function' || (typeof XMLHttpRequest === "undefined" ? "undefined" : _typeof(XMLHttpRequest)) === 'object';
}
function isPromise(maybePromise) {
  return !!maybePromise && typeof maybePromise.then === 'function';
}
export function makePromise(maybePromise) {
  if (isPromise(maybePromise)) {
    return maybePromise;
  }
  return Promise.resolve(maybePromise);
}
var interpolationRegexp = /\{\{(.+?)\}\}/g;
export function interpolate(str, data) {
  return str.replace(interpolationRegexp, function (match, key) {
    var k = key.trim();
    if (UNSAFE_KEYS.indexOf(k) > -1) return match;
    var value = data[k];
    return value != null ? value : match;
  });
}
export function interpolateUrl(str, data) {
  var unsafe = false;
  var result = str.replace(interpolationRegexp, function (match, key) {
    var k = key.trim();
    if (UNSAFE_KEYS.indexOf(k) > -1) return match;
    var value = data[k];
    if (value == null) return match;
    var check = SAFETY_CHECK_BY_KEY[k] || isSafeLangUrlSegment;
    var segments = String(value).split('+');
    var _iterator = _createForOfIteratorHelper(segments),
      _step;
    try {
      for (_iterator.s(); !(_step = _iterator.n()).done;) {
        var seg = _step.value;
        if (!check(seg)) {
          unsafe = true;
          return match;
        }
      }
    } catch (err) {
      _iterator.e(err);
    } finally {
      _iterator.f();
    }
    return segments.join('+');
  });
  return unsafe ? null : result;
}