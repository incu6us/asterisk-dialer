export const DROPDOWN_SELECTOR = '[data-dropdown]';
export const DROPDOWN_ATTR = 'data-dropdown';
export const TRIGGER_CLASS_NAME = 'dropdown-trigger';
export const CONTENT_CLASS_NAME = 'dropdown-content';
const CONTENT_INDENT = 5;

/**
 * Return scroll top position from top of the document.
 * @returns {Number}
 */
export function scrollTopOfDocument() {
    return window.scrollY
        || window.pageYOffset
        || document.documentElement.scrollTop;
}

/**
 * @returns {Number}
 */
export function scrollLeftOfDocument() {
    return window.scrollX
        || window.pageXOffset
        || document.documentElement.scrollLeft;
}

/**
 * @param {Element} triggerEl
 * @param {Element} contentEl
 * @param {Event} event
 */
function handleWindowClick (triggerEl, contentEl, event) {
    if (!document.body.contains(event.target) || !event.isTrusted) {
        event.stopPropagation();
        event.preventDefault();
        return;
    }
    let node = event.target;
    while (node.parentNode) {
        if (node.classList.contains('dropdown-content')) {
            return;
        }
        node = node.parentNode;
    }
    if (!triggerEl.contains(event.target) && !contentEl.contains(event.target)) {
        hide(triggerEl, contentEl);
    }
}

/**
 * @param {Element} triggerEl
 * @param {Element} contentEl
 * @param {Event} event
 */
function handleWindowResize (triggerEl, contentEl, event) {
    updateContentPosition(triggerEl, contentEl);
}

/**
 * @param {Element} contentEl
 */
function updateContent (contentEl) {
    contentEl.classList.add(CONTENT_CLASS_NAME);
    if (contentEl.parentNode!=document.body) {
        Array.prototype.forEach.call(document.querySelectorAll('[id="'+contentEl.getAttribute('id')+'"]'), (el) => {
            if (el!==contentEl) {
                el.parentNode.removeChild(el);
            }
        });
        contentEl.parentNode.removeChild(contentEl);
        document.body.appendChild(contentEl);
    }
}

/**
 * @param {Element} triggerEl
 */
function updateTrigger (triggerEl) {
    triggerEl.classList.add(TRIGGER_CLASS_NAME);
    triggerEl.addEventListener('click', handleTriggerClick);
}

/**
 * @param {Event} event
 */
function handleTriggerClick (event) {
    event.preventDefault();
    const triggerEl = event.currentTarget;
    const contentEl = findContentByTrigger(triggerEl);
    const isClosableOnTriggerClick = triggerEl.getAttribute('data-is-closable-on-trigger-click')==='1';
    if (contentEl.classList.contains('is-visible') && isClosableOnTriggerClick) {
        hide(triggerEl, contentEl);
    } else {
        show(triggerEl, contentEl);
    }
}

/**
 * @param {Element} triggerEl
 * @param {Element} contentEl
 */
function updateContentPosition (triggerEl, contentEl) {
    const scrollTop = scrollTopOfDocument();
    const scrollLeft = scrollLeftOfDocument();
    const boundingRectOfContent = contentEl.getBoundingClientRect();
    const boundingRectOfTrigger = triggerEl.getBoundingClientRect();
    const bottom = window.innerHeight-(boundingRectOfTrigger.top+boundingRectOfContent.height);
    let verticalPosition = null;
    let horizontalPosition = null;
    if (bottom>boundingRectOfContent.height || boundingRectOfTrigger.top<boundingRectOfContent.height) {
        contentEl.style.top = (scrollTop+boundingRectOfTrigger.top+boundingRectOfTrigger.height+CONTENT_INDENT)+'px';
        verticalPosition = 'bottom';
    } else {
        contentEl.style.top = (scrollTop+boundingRectOfTrigger.top-boundingRectOfContent.height-CONTENT_INDENT)+'px';
        verticalPosition = 'top';
    }
    if (boundingRectOfTrigger.left>boundingRectOfContent.width) {
        contentEl.style.left = (scrollLeft+boundingRectOfTrigger.left-boundingRectOfContent.width+boundingRectOfTrigger.width)+'px';
        horizontalPosition = 'right';
    } else {
        contentEl.style.left = (scrollLeft+boundingRectOfTrigger.left)+'px';
        horizontalPosition = 'left';
    }
    contentEl.className = contentEl.className.replace(new RegExp(`${ CONTENT_CLASS_NAME }-(top|bottom)-(right|left)`, 'ig'), '');
    contentEl.classList.add(`${ CONTENT_CLASS_NAME }-${ verticalPosition }-${ horizontalPosition }`);
}

/**
 * @param {String} id
 * @returns {Element}
 */
function findContent (id) {
    return document.getElementById(id);
}

/**
 * @param {Element} triggerEl
 * @returns {Element}
 */
function findContentByTrigger (triggerEl) {
    return findContent(triggerEl.getAttribute(DROPDOWN_ATTR));
}

/**
 * Show Dropdown.
 * @param {Element} triggerEl
 * @param {Element} contentEl
 */
export function show (triggerEl, contentEl) {
    if (!contentEl.getAttribute('data-is-resize-click-handlers-attached')) {
        window.addEventListener('resize', handleWindowResize.bind(null, triggerEl, contentEl), false);
        document.body.addEventListener('click', handleWindowClick.bind(null, triggerEl, contentEl), false);
        contentEl.setAttribute('data-is-resize-click-handlers-attached', '1');
    }
    triggerEl.classList.add('is-opened');
    contentEl.classList.add(CONTENT_CLASS_NAME);
    contentEl.classList.add('is-visible');
    updateContentPosition(triggerEl, contentEl);
}

/**
 * Hide Dropdown.
 * @param {Element} triggerEl
 * @param {Element} contentEl
 */
export function hide (triggerEl, contentEl) {
    window.removeEventListener('resize', handleWindowResize.bind(null, triggerEl, contentEl), false);
    document.body.removeEventListener('click', handleWindowClick.bind(null, triggerEl, contentEl), false);
    contentEl.classList.remove('is-visible');
    const event = document.createEvent('CustomEvent');
    event.initCustomEvent('hide', true, true, {});
    contentEl.dispatchEvent(event);
}

/**
 * @param {Element} triggerEl
 */
export function hideByTrigger (triggerEl) {
    hide(triggerEl, findContentByTrigger(triggerEl));
}

/**
 * Update Dropdown trigger and content elements.
 * @param {String} selector
 * @param {Element} scope
 */
export function updateBySelector (selector=null, scope=null, openedIds=[]) {
    selector = selector || DROPDOWN_SELECTOR;
    scope = scope || document;
    Array.prototype.forEach.call(scope.querySelectorAll(selector), (triggerEl) => {
        const contentEl = findContentByTrigger(triggerEl);
        updateTrigger(triggerEl);
        updateContent(contentEl);
        if (openedIds.indexOf(contentEl.getAttribute('id')) !== -1) {
            show(triggerEl, contentEl);
        } else {
            hide(triggerEl, contentEl);
        }
    });
}

/**
 * @param {Element} triggerEl
 * @param {Element} contentEl
 * @returns {boolean}
 */
export function isActive (triggerEl, contentEl) {
    return contentEl.classList.contains('is-visible');
}
