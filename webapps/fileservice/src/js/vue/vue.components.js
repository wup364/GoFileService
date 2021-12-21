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
 *  自动适应高度, 有两种传值方式: number | function
 *  传入function(height): 通过一个回调函数来告诉调用者当前父级的高度, 再在函数里面去刷新高度, 这种方式没有警告. 如: v-minus-height="(h) => (tableHeight = h - 130)"
 *  传入number: 传入一个数值时, 程序自动改变当前元素告诉, 但控制台可能会有[Vue warring]. 虽不影响功能,但推荐使用函数方式. 如: v-minus-height="130"
 */
Vue.directive('minus-height', {
	// 绑定钩子函数
	bind(el, binding, vnode) { },
	// 绑定到节点函数
	inserted: () => { },
	// 组件更新钩子函数
	update(el, binding, vnode) { },
	// 组件更新完成
	componentUpdated(el, binding, vnode, vnodeold) {
		if (vnode.v_unbind_minus_height) {
			vnode.v_unbind_minus_height();
		}
		if (vnodeold.v_unbind_minus_height) {
			vnodeold.v_unbind_minus_height();
		}
		// 高度自动处理函数
		vnode.v_minus_height = (listen) => {
			try {
				if (el.parentNode.clientHeight == 0) {
					return;
				}
				if (typeof binding.value === 'function') {
					// 通过一个回调函数来告诉调用者当前父级的高度, 这种方式没有警告
					binding.value(el.parentNode.clientHeight);
				} else {
					// 当绑定对象不是一个函数时, 控制台可能会有 [Vue warring], 但不影响功能. 推荐使用函数方式
					vnode.componentInstance.height = el.parentNode.clientHeight - binding.value;
				}
			} catch (e) { }
			if (listen === true) {
				window.addEventListener("resize", vnode.v_minus_height, false);
			}
		};
		// 接触高度变化监控
		vnode.v_unbind_minus_height = () => {
			window.removeEventListener("resize", vnode.v_minus_height, false);
			vnode.v_minus_height = undefined;
			vnode.v_unbind_minus_height = undefined;
		}
		// 监听导读变化
		vnode.v_minus_height(true);
	},
	// 解除指令 
	unbind(el, binding, vnode) {
		if (vnode.v_unbind_minus_height) {
			vnode.v_unbind_minus_height();
		}
	}
});
// 
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