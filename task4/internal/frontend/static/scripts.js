// Global state
let currentData = null;

// Load data from API
async function loadData(type) {
    try {
        showLoading();

        const response = await fetch(`/api/${type}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        currentData = data;
        displayResults(data);
        updateTimestamp();

    } catch (error) {
        console.error('Error loading data:', error);
        showError('Failed to load data: ' + error.message);
    }
}

// Allocate memory for testing
async function allocateMemory(size) {
    try {
        const response = await fetch(`/api/alloc?size=${size}`);
        const result = await response.json();

        showNotification(`Successfully allocated ${result.mb.toFixed(2)} MB`);

        // Auto-refresh system stats
        setTimeout(() => loadData('system'), 500);

    } catch (error) {
        console.error('Error allocating memory:', error);
        showError('Failed to allocate memory');
    }
}

// Display results in the UI
function displayResults(data) {
    const resultsDiv = document.getElementById('results');

    if (data.type === 'system') {
        resultsDiv.innerHTML = createSystemView(data);
    } else {
        resultsDiv.innerHTML = createSingleView(data);
    }
}

// Create view for system data (multiple sections)
function createSystemView(data) {
    return `
        <div class="report-section">
            <div class="section-header-report">
                <span>⚡</span>
                <h3>System Overview</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">${formatJSON(data.data)}</div>
            </div>
        </div>
        
        <div class="report-section">
            <div class="section-header-report">
                <span>📝</span>
                <h3>Allocations</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">${formatJSON(data.data.Allocations)}</div>
            </div>
        </div>
        
        <div class="report-section">
            <div class="section-header-report">
                <span>🔄</span>
                <h3>Garbage Collector</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">${formatJSON(data.data.GC)}</div>
            </div>
        </div>
        
        <div class="report-section">
            <div class="section-header-report">
                <span>💾</span>
                <h3>Memory Usage</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">${formatJSON(data.data.Memory)}</div>
            </div>
        </div>
        
        <div class="report-section">
            <div class="section-header-report">
                <span>📊</span>
                <h3>Runtime Info</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">Goroutines: ${data.data.Goroutines}\nTimestamp: ${new Date(data.data.Timestamp).toLocaleString()}</div>
            </div>
        </div>
    `;
}

// Create view for single data type
function createSingleView(data) {
    return `
        <div class="report-section">
            <div class="section-header-report">
                <span>${getIcon(data.type)}</span>
                <h3>${capitalizeFirst(data.type)} Report</h3>
            </div>
            <div class="section-content-report">
                <div class="json-data">${formatJSON(data.data)}</div>
            </div>
        </div>
    `;
}

// Format JSON for display
function formatJSON(obj) {
    return JSON.stringify(obj, null, 2)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/\n/g, '<br>')
        .replace(/ /g, '&nbsp;');
}

// Utility functions
function capitalizeFirst(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

function getIcon(type) {
    const icons = {
        'system': '⚡',
        'allocations': '📝',
        'gc': '🔄',
        'memory': '💾'
    };
    return icons[type] || '📊';
}

function showLoading() {
    document.getElementById('results').innerHTML = `
        <div class="loading">Loading data</div>
    `;
}

function showError(message) {
    document.getElementById('results').innerHTML = `
        <div class="empty-state">
            <div class="empty-icon">❌</div>
            <h3>Error</h3>
            <p>${message}</p>
        </div>
    `;
}

function showNotification(message) {
    // Simple notification
    alert(message);
}

function updateTimestamp() {
    document.getElementById('timestamp').textContent =
        `Last updated: ${new Date().toLocaleString()}`;
}

function clearResults() {
    document.getElementById('results').innerHTML = `
        <div class="empty-state">
            <div class="empty-icon">📊</div>
            <h3>No data loaded</h3>
            <p>Click on any button above to load metrics data</p>
        </div>
    `;
    document.getElementById('timestamp').textContent = '';
    currentData = null;
}

function copyToClipboard() {
    if (!currentData) {
        showNotification('No data to copy');
        return;
    }

    const text = JSON.stringify(currentData, null, 2);
    navigator.clipboard.writeText(text).then(() => {
        showNotification('Data copied to clipboard!');
    }).catch(err => {
        console.error('Failed to copy: ', err);
        showNotification('Failed to copy data');
    });
}

function openPrometheusPage() {
    window.open('/metrics', '_blank');
}

function refreshAll() {
    loadData('system');
}

// Initialize
document.addEventListener('DOMContentLoaded', function () {
    // Load system data on page load
    loadData('system');
});