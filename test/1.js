window.joiner = function(sep) {
  this.sep = sep;
}
window.joiner.prototype = {
  join: function(ary) {
    return ary.join(this.sep)
  }
}
