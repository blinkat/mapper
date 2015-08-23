// mapper havs 3 layer
// layer#1 = the map graph
// layer#2 = mask
// layer#3 = label
// layer#3 = info win

+ function() {
    // global config
    var $doc = $(document);

    $.mapper = {
        handler: "map/handler",
        maxMove: 64,
    }

    // grid object
    function Grid(father) {
        var $this = this;
        this.size = father.mapper.size;
        this.element = $(document.createElement("div")).addClass("grid").width(this.size).height(this.size);
        this.$img = $(document.createElement('img')).attr('alt', '').width(this.size).height(this.size)
        this.mapCoord = {
            x: 0,
            y: 0,
        }

        var $this = this;


        // private methods
        function init() {
            var ele = $this.element;
            father.element.append(ele);
            ele.append($this.$img)
        }

        // public methods
        this.top = function(val) {
            return val != null ? $this.element.css('top', val + 'px') : getCssNum($this.element.css('top'));
        }

        this.left = function(val) {
            return val != null ? $this.element.css('left', val + 'px') : getCssNum($this.element.css('left'));
        }

        this.image = function(x, y) {
            // $this.$img.attr('src', '')
            $this.x = x
            $this.y = y
            $this.$img.attr('src', $.mapper.handler + "?type=get&x=" + x + "&y=" + y + "&scale=" + father.mapper.scaleM)
        }

        this.x = 0
        this.y = 0

        init();
    }

    function ImgLayer(father) {
        this.grids = []
        this.element = $(document.createElement("div")).addClass("mapper-img")
        this.mapper = father

        var $this = this,
            offset = {}

        function init() {
            var wide = Math.floor(father.element.width() / 2)
            var high = Math.floor(father.element.height() / 2)
            $this.left(wide)
            $this.top(high)

            father.element.append($this.element)
            $this.resize()

            offset.left = ((wide >> 8) + 1) * father.size
            offset.top = ((high >> 8) + 1) * father.size
        }

        this.resize = function() {
            var size = {
                    x: Math.floor(father.element.width() / father.size) + 1,
                    y: Math.floor(father.element.height() / father.size) + 1,
                },
                gridsize = father.size;

            size.x = size.x % 2 ? size.x + 1 : size.x;
            size.y = size.y % 2 ? size.y + 1 : size.y;
            size.x = (size.x / 2) + 1
            size.y = (size.y / 2) + 1

            for (var x = -size.x; x < size.x; x++) {
                var line = $this.grids[x + size.x]
                if (!line) $this.grids.push(line = [])
                for (var y = -size.y; y < size.y; y++) {
                    var g = line[y + size.y]
                    if (!g) {
                        g = new Grid($this)
                        g.left(x * gridsize)
                        g.top(y * gridsize)
                        g.x = x
                        g.y = y
                        line.push(g)
                    }
                    //g.image(g.x, g.y)
                }
            }
        }

        this.draw = function() {
            for (var x in $this.grids) {
                var line = $this.grids[x]
                for (var y in line) {
                    var g = line[y]
                    g.image(g.x, g.y)
                }
            }
        }

        this.top = function(val) {
            return val != null ? $this.element.css('top', val + 'px') : getCssNum($this.element.css('top'));
        }

        this.left = function(val) {
            return val != null ? $this.element.css('left', val + 'px') : getCssNum($this.element.css('left'));
        }

        this.move = function(x, y) {
            if (x >= father.size) {
                console.log(x)
            }
            x = x >= father.size ? $.mapper.maxMove : x
            y = y >= father.size ? $.mapper.maxMove : y
            var pl = $this.left(),
                pt = $this.top(),
                maxl = father.center.x,
                maxt = father.center.y,
                tmp = maxl

            if (pl + x > tmp) {
                if (pl == tmp) {
                    x = 0
                    father.mask.warring("left")
                } else x = tmp - pl
            } else {
                tmp = -tmp + father.element.width()
                if (pl + x < tmp) {
                    if (pl == tmp) {
                        x = 0
                        father.mask.warring("right")
                    } else x = tmp - pl
                }
            }

            tmp = maxt
            if (pt + y > tmp) {
                if (pt == tmp) {
                    father.mask.warring("top")
                    y = 0
                } else y = tmp - pt
            } else {
                tmp = -tmp + father.element.height()
                if (pt + y < tmp) {
                    if (pt == tmp) {
                        y = 0
                        father.mask.warring("bottom")
                    } else y = tmp - pt
                }
            }

            // set

            if (x != 0) {
                pl += x
                $this.left(pl)

                var offset_left = offset.left - pl
                if (offset_left <= 0) {
                    var first = $this.grids[0]
                    var last = $this.grids.pop()
                    var left = first[0].left() - father.size
                    var x = first[0].x - 1
                    for (var i = 0; i < last.length; i++) {
                        var grid = last[i]
                        grid.left(left)
                        grid.image(x, grid.y)
                    }
                    $this.grids.unshift(last)
                    offset.left += father.size
                } else if (offset_left >= father.size) {
                    var first = $this.grids.shift()
                    var last = $this.grids[$this.grids.length - 1]
                    var left = last[0].left() + father.size
                    var x = last[0].x + 1
                    for (var i = 0; i < first.length; i++) {
                        first[i].left(left)
                        first[i].image(x, first[i].y)
                    }
                    $this.grids.push(first)
                    offset.left -= father.size
                }
            }

            if (y != 0) {
                pt += y
                $this.top(pt)
                var offset_top = offset.top - pt

                if (offset_top <= 0) {
                    var leng = $this.grids.length
                    var cha = $this.grids[0][0].top() - father.size
                    var y = $this.grids[0][0].y - 1
                    for (var i = 0; i < leng; i++) {
                        var a = $this.grids[i]
                        var bot = a.pop()
                        bot.top(cha)
                        bot.image(bot.x, y)
                        a.unshift(bot)
                    }
                    offset.top += father.size
                } else if (offset_top >= father.size) {
                    var leng = $this.grids.length
                    var leng2 = $this.grids[0].length - 1
                    var y = $this.grids[0][leng2].y + 1
                    var cha = $this.grids[0][leng2].top() + father.size
                    for (var i = 0; i < leng; i++) {
                        var a = $this.grids[i]
                        var top = a.shift()
                        top.top(cha)
                        top.image(top.x, y)
                        a.push(top)
                    }
                    offset.top -= father.size
                }
            }
        }

        init();
    }

    function MaskLayer(father) {
        var $this = this;
        this.element = $(document.createElement("div")).addClass('mapper-mask').attr("oncontextmenu", "return false")
        this.element.attr("ondragstart", "return false").attr("onselectstart", "return false").attr("onselect", "return false")
        var isdown = false,
            prevcoord = {},
            moveTimer,
            prevbuffer;

        var warring = {}


        function init() {
            father.element.append($this.element)
            warring.top = $(document.createElement('div')).css('background-image',
                'url(' + $.mapper.handler + "?type=warring-img&img-type=top)").addClass("mapper-mask-horizonal-warring mapper-top")
            warring.left = $(document.createElement('div')).css('background-image',
                'url(' + $.mapper.handler + "?type=warring-img&img-type=left)").addClass("mapper-mask-verical-warring mapper-left")
            warring.bottom = $(document.createElement('div')).css('background-image',
                'url(' + $.mapper.handler + "?type=warring-img&img-type=bottom)").addClass("mapper-mask-horizonal-warring mapper-bottom")
            warring.right = $(document.createElement('div')).css('background-image',
                'url(' + $.mapper.handler + "?type=warring-img&img-type=right)").addClass("mapper-mask-verical-warring mapper-right")

            $this.element.append(warring.top)
            $this.element.append(warring.left)
            $this.element.append(warring.bottom)
            $this.element.append(warring.right)
        }

        function getmovedis(event) {
            return {
                x: event.clientX - prevcoord.x,
                y: event.clientY - prevcoord.y,
            }
        }

        function run_animate(v) {
            if (v.war_timer) {
                clearTimeout(v.war_timer)
            }
            v.addClass("warring")
            v.war_timer = setTimeout(function() {
                v.removeClass("warring")
            }, 400)
        }

        this.warring = function(v) {
            run_animate(warring[v])
        }

        // events
        this.element.mousedown(function(event) {
            isdown = true
            if (moveTimer) {
                clearTimeout(moveTimer)
            }

            prevcoord = {
                x: event.clientX,
                y: event.clientY,
            }
        });

        $doc.mousemove(function(event) {
            if (isdown) {
                clearTimeout(moveTimer)
                var coord = getmovedis(event)
                prevcoord.x = event.clientX
                prevcoord.y = event.clientY

                father.img.move(coord.x, coord.y)

                prevbuffer = coord
            }
        });

        $doc.mouseup(function(event) {
            if (isdown) {
                isdown = false
                var coord = prevbuffer

                if (Math.abs(coord.x) <= 2 && Math.abs(coord.y) <= 2) {
                    return
                }

                var modulus = 1.5;

                function timer() {
                    var dis = {
                        x: Math.floor(coord.x * modulus),
                        y: Math.floor(coord.y * modulus),
                    }
                    modulus -= 0.05
                    father.img.move(dis.x, dis.y)
                    if (modulus > 0) {
                        moveTimer = setTimeout(timer, 25)
                    }
                }

                timer()
            }
        });

        $doc.mousewheel(function (event) {
            var y = event.deltaY
            if ( y > 0) {
                father.scale(father.scaleM + 1)
            } else if (y < 0) {
                father.scale(father.scaleM - 1)
            }
        })

        init();
    }

    // mapper object
    // @element = map html object
    // @opt = all option params
    function Mapper(element, opt) {
        var map = this;
        this.element = element;
        this.option = opt;
        this.center = {}
            // girds
        this.grids = [];
        this.scaleM = opt.scale

        // private methods
        // initizle map contents
        function init() {
            $.getJSON($.mapper.handler, "type=init-data", function(res) {
                if (res && typeof res == "object") {
                    map.size = res.size
                    map.element.addClass('mapper');
                    if ($.map.Pointer) map.element.css('cursor', 'url(' + $.map.Pointer + ")");

                    init_layer();
                    map.scale(map.option.scale)
                } else {
                    throw "get init data faild.";
                }
            });
        }

        function init_layer() {
            map.img = new ImgLayer(map);
            map.mask = new MaskLayer(map);
        }

        // public methods

        // move map, x,y = offset value
        this.move = function(x, y) {
            var pos = {
                y: getCssNum(layer.css('margin-top')),
                x: getCssNum(layer.css('margin-left')),
            }
        }

        this.scale = function(scale) {
            $.getJSON($.mapper.handler, "type=cent&scale=" + scale, function(res) {
                if (res && typeof res == "object") {
                    map.scaleM = scale
                    map.center = res
                    map.img.draw()
                }
            })
        }

        init(); // goto init
    }

    Mapper._DEFAULT = {
        scale: 4,
    };

    Mapper.VERSION = '0.0.1';

    // plugin
    $.fn.mapper = function(opt) {
        var $this = $(this);
        var option = typeof opt == "object" && opt
        var data = $this.data('mapper');

        if (!data) {
            $this.data('mapper', (data = new Mapper($this, $.extend({}, Mapper._DEFAULT, option))));
        }
        if (option != null && typeof option == 'string') {
            return data.do(option) || data;
        }
        return data;
    }

    // helper func
    function getCssNum(css) {
        if (css == "auto") return 0
        var rege = /^[0-9]*$/;
        var i = css.length;
        for (; i >= 0; i--) {
            if (rege.test(css.charAt(i))) {
                break;
            }
        }

        css = css.substr(0, i);
        return parseInt(css);
    }
}(jQuery)