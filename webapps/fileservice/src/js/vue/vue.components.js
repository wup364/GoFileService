// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE组件, 无业务属性
"use strict";
import Vue from 'vue'
/**
 *  监听高度变化, 入参为 function(当前节点高度, 父级节点高度)
 *  通过一个回调函数告诉调用者当前的高度. 如: v-watch-height="(ch,ph) => (xxxHeight = ph - 130)"
 */
Vue.directive('watch-height', {
	// 绑定钩子函数
	bind(el, binding, vnode) { },
	// 绑定到节点函数
	inserted: () => { },
	// 组件更新钩子函数
	update(el, binding, vnode) { },
	// 组件更新完成
	componentUpdated(el, binding, vnode, vnodeold) {
		if (vnode.v_unbind_watch_height) {
			vnode.v_unbind_watch_height();
		}
		if (vnodeold.v_unbind_watch_height) {
			vnodeold.v_unbind_watch_height();
		}
		if (typeof binding.value !== 'function') {
			return
		}
		// 高度自动处理函数
		vnode.v_watch_height = (listen) => {
			try {
				binding.value(el.clientHeight, el.parentNode.clientHeight);
			} catch (e) { console.error(err); }
			if (listen === true) {
				window.addEventListener("resize", vnode.v_watch_height, false);
			}
		};
		// 接触高度变化监控
		vnode.v_unbind_watch_height = () => {
			window.removeEventListener("resize", vnode.v_watch_height, false);
			vnode.v_watch_height = undefined;
			vnode.v_unbind_watch_height = undefined;
		}
		// 监听导读变化
		vnode.v_watch_height(true);
	},
	// 解除指令 
	unbind(el, binding, vnode) {
		if (vnode.v_unbind_watch_height) {
			vnode.v_unbind_watch_height();
		}
	}
});

/**
 * 右键菜单
 * menus = {id: {name: 'xxx', icon:'', show:true, handler: func, divided:false, }}
 */
Vue.component("right-click-menu", {
	props: ['bindRef', 'menus'],
	data() {
		return {
			posX: 0,
			posY: 0,
			currentVisible: false
		}
	},
	template: `
	<Dropdown :style='locatorStyle' placement='right-start' trigger='custom' :visible='currentVisible' @on-click='onClick' @on-clickoutside='handleCancel'>
		<Dropdown-Menu slot='list'>
			<Dropdown-Item v-for='(val, key) in menus' v-show='val.show' :name='key' :divided='val.divided'>
			<i v-if='val.icon' :class='val.icon' style='padding-right: 3px;'></i> <span style='vertical-align: middle;'>{{val.name}}</span>
			</Dropdown-Item>
		</Dropdown-Menu>
	</Dropdown>
	`,
	computed: {
		locatorStyle() {
			return {
				position: 'fixed',
				left: this.posX + 'px',
				top: this.posY + 'px',
				maxHeight: 'unset',
				binddom: null,
			}
		}
	},
	methods: {
		onClick(name) {
			this.currentVisible = false
			if (this.menus[name] && this.menus[name].handler) {
				try {
					this.menus[name].handler(this.menus[name]);
				} catch (e) { console.error(e); }
			}
		},
		handleContextmenu(e) {
			e.preventDefault();
			e.stopPropagation();
			e.returnValue = false;
			if (e.button === 2) {
				this.currentVisible = false;
				if (this.posX !== e.clientX) { this.posX = e.clientX; }
				if (this.posY !== e.clientY) { this.posY = e.clientY; }
				this.$nextTick(() => {
					this.currentVisible = true;
				});
			}
		},
		handleCancel() {
			this.currentVisible = false
		},
		getRefNode() {
			if (!this.bindRef) {
				return document;
			}
			let node = this;
			while (true) {
				if (node.$refs[this.bindRef]) {
					break;
				}
				node = node.$parent;
				if (!node) { break; }
			}
			return node ? node.$refs[this.bindRef].$el : {};
		},
	},
	mounted() {
		this.binddom = this.getRefNode();
		this.binddom.addEventListener('contextmenu', this.handleContextmenu, true);
		this.binddom.addEventListener('mouseup', this.handleCancel, true);
	},
	destroyed() {
		try {
			this.currentVisible = false;
			this.binddom.removeEventListener('contextmenu', this.handleContextmenu, true);
			this.binddom.removeEventListener('mouseup', this.handleCancel, true);
		} catch (e) { }
	},
	watch: {
	}
});