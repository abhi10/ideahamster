// Idea Hamster - Main Application JavaScript
(function() {
  'use strict';

  // ============================================
  // HTMX CSRF Configuration
  // ============================================
  function initCSRF() {
    document.body.addEventListener('htmx:configRequest', function(event) {
      var csrfToken = document.querySelector('meta[name="csrf-token"]');
      if (csrfToken) {
        event.detail.headers['X-CSRF-Token'] = csrfToken.content;
      }
    });
  }

  // ============================================
  // Modal Handlers
  // ============================================
  function showModal(modalId) {
    var modal = document.getElementById(modalId);
    if (modal) {
      modal.style.display = 'flex';
    }
  }

  function hideModal(modalId) {
    var modal = document.getElementById(modalId);
    if (modal) {
      modal.style.display = 'none';
    }
  }

  // Close modal when clicking outside content
  function initModalBackdropClose() {
    document.addEventListener('click', function(event) {
      if (event.target.classList.contains('modal-backdrop')) {
        event.target.style.display = 'none';
      }
    });
  }

  // Close modal on Escape key
  function initModalEscapeClose() {
    document.addEventListener('keydown', function(event) {
      if (event.key === 'Escape') {
        var modals = document.querySelectorAll('[id$="-modal"]');
        modals.forEach(function(modal) {
          if (modal.style.display === 'flex') {
            modal.style.display = 'none';
          }
        });
      }
    });
  }

  // ============================================
  // Event Delegation for Dynamic Content
  // ============================================
  function initEventDelegation() {
    document.addEventListener('click', function(event) {
      var target = event.target;

      // Handle modal close buttons
      if (target.hasAttribute('data-close-modal')) {
        var modalId = target.getAttribute('data-close-modal');
        hideModal(modalId);
        return;
      }

      // Handle modal open buttons
      if (target.hasAttribute('data-open-modal')) {
        var modalId = target.getAttribute('data-open-modal');
        showModal(modalId);
        return;
      }

      // Handle cancel buttons inside modals (traverse up to find button)
      var cancelBtn = target.closest('[data-action="cancel"]');
      if (cancelBtn) {
        var modal = cancelBtn.closest('[id$="-modal"]');
        if (modal) {
          modal.style.display = 'none';
        }
        return;
      }
    });
  }

  // ============================================
  // HTMX Custom Event Handlers
  // ============================================
  function initHTMXEvents() {
    // Listen for showModal trigger from server
    document.body.addEventListener('showModal', function(event) {
      var modalId = event.detail.value;
      showModal(modalId);
    });

    // Listen for closeModal trigger from server
    document.body.addEventListener('closeModal', function(event) {
      var modalId = event.detail.value;
      hideModal(modalId);
    });

    // Listen for showNotification trigger from server
    document.body.addEventListener('showNotification', function(event) {
      var message = event.detail.value;
      // Simple notification - could be enhanced with a toast library
      if (message) {
        console.log('Notification:', message);
      }
    });
  }

  // ============================================
  // Initialization
  // ============================================
  function init() {
    initCSRF();
    initHTMXEvents();
    initEventDelegation();
    initModalBackdropClose();
    initModalEscapeClose();
  }

  // Expose functions globally for any inline needs during transition
  window.IdeaHamster = {
    showModal: showModal,
    hideModal: hideModal,
    init: init
  };

  // Auto-initialize on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
