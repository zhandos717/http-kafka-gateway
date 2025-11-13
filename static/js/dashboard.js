// Configuration
const config = {
    refreshInterval: 3000, // 3 seconds
    metricsEndpoint: '/api/metrics',
    topicsEndpoint: '/api/topics',
    messagesEndpoint: '/api/messages',
    sendMessageEndpoint: '/message',
    healthEndpoint: '/health'
};

// Global state
let metrics = {
    messagesProcessed: { success: 0, error: 0 },
    requestDuration: 0,
    kafkaErrors: 0,
    authAttempts: { success: 0, error: 0 },
    httpLatency: 0
};

// DOM Elements
const elements = {
    totalMessages: document.getElementById('total-messages'),
    successfulMessages: document.getElementById('successful-messages'),
    errorMessages: document.getElementById('error-messages'),
    avgResponseTime: document.getElementById('avg-response-time'),
    messageLog: document.getElementById('message-log'),
    topicsList: document.getElementById('topics-list'),
    statusIndicator: document.getElementById('status-indicator'),
    statusText: document.getElementById('status-text'),
    messageForm: document.getElementById('message-form'),
    valueTextarea: document.getElementById('value'),
    formatJsonBtn: document.getElementById('format-json-btn')
};

// Initialize the dashboard
async function initDashboard() {
    await fetchMetrics();
    setupEventListeners();
    
    // Start periodic updates
    setInterval(fetchMetrics, config.refreshInterval);
}

// Format JSON in the value textarea
function formatJson() {
    const valueTextarea = elements.valueTextarea;
    const currentValue = valueTextarea.value.trim();
    
    if (!currentValue) {
        return;
    }
    
    try {
        const parsedJson = JSON.parse(currentValue);
        const formattedJson = JSON.stringify(parsedJson, null, 2);
        valueTextarea.value = formattedJson;
    } catch (error) {
        alert('Invalid JSON: ' + error.message);
    }
}

// Set up event listeners
function setupEventListeners() {
    elements.messageForm.addEventListener('submit', handleSendMessage);
    elements.formatJsonBtn.addEventListener('click', formatJson);
}

// Fetch metrics from the server
async function fetchMetrics() {
    try {
        const response = await fetch('/api/metrics');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        const metrics = data.metrics || {};
        
        // Format the metrics for the dashboard
        const formattedMetrics = {
            messagesProcessed: {
                success: metrics.messages_processed?.success || 0,
                error: metrics.messages_processed?.error || 0
            },
            kafkaErrors: metrics.kafka_errors || 0,
            authAttempts: {
                success: metrics.auth_attempts?.success || 0,
                error: metrics.auth_attempts?.error || 0
            },
            requestDuration: (metrics.request_duration_sum || 0).toFixed(2)
        };

        updateMetrics(formattedMetrics);
        updateStatusIndicator(true);
        
        // Update logs with recent activity
        await updateLogs();
        
    } catch (error) {
        console.error('Error fetching metrics:', error);
        updateStatusIndicator(false);
    }
}

// Update metrics display
function updateMetrics(newMetrics) {
    metrics = newMetrics;
    
    elements.totalMessages.textContent = (metrics.messagesProcessed.success + metrics.messagesProcessed.error).toLocaleString();
    elements.successfulMessages.textContent = metrics.messagesProcessed.success.toLocaleString();
    elements.errorMessages.textContent = metrics.messagesProcessed.error.toLocaleString();
    elements.avgResponseTime.textContent = `${metrics.requestDuration}ms`;
}

// Update status indicator
function updateStatusIndicator(isConnected) {
    if (isConnected) {
        elements.statusIndicator.className = 'w-3 h-3 rounded-full bg-green-500 mr-2';
        elements.statusText.textContent = 'Connected';
        elements.statusText.className = 'text-sm font-medium text-green-600';
    } else {
        elements.statusIndicator.className = 'w-3 h-3 rounded-full bg-red-500 mr-2';
        elements.statusText.textContent = 'Disconnected';
        elements.statusText.className = 'text-sm font-medium text-red-600';
    }
}

// Update logs with recent activity
async function updateLogs() {
    try {
        // Fetch recent messages
        const messagesResponse = await fetch(config.messagesEndpoint);
        if (messagesResponse.ok) {
            const messagesData = await messagesResponse.json();
            const messages = messagesData.messages || [];
            
            // Clear existing message logs
            elements.messageLog.innerHTML = '';
            
            // Populate with recent messages
            messages.slice(0, 5).forEach(message => {
                const messageItem = document.createElement('li');
                messageItem.className = 'py-2';
                const timestamp = new Date(message.timestamp).toLocaleTimeString();
                
                messageItem.innerHTML = `
                    <div class="flex justify-between">
                        <div class="text-sm">
                            <p class="font-medium text-gray-900">${message.topic}</p>
                            <p class="text-gray-500">Message</p>
                        </div>
                        <div class="text-sm">
                            <p class="font-medium ${message.status === 'success' ? 'text-green-600' : 'text-red-600'}">
                                ${message.status.toUpperCase()}
                            </p>
                            <p class="text-gray-500">${timestamp}</p>
                        </div>
                    </div>
                `;
                
                elements.messageLog.appendChild(messageItem);
            });
        }
    } catch (error) {
        console.error('Error fetching recent messages:', error);
    }
    
    try {
        // Fetch topics
        const topicsResponse = await fetch(config.topicsEndpoint);
        if (topicsResponse.ok) {
            const topicsData = await topicsResponse.json();
            const topics = topicsData.topics || [];
            
            // Clear existing topics list
            elements.topicsList.innerHTML = '';
            
            // Populate with topics
            topics.slice(0, 5).forEach(topic => {
                const topicItem = document.createElement('li');
                topicItem.className = 'py-2';
                
                topicItem.innerHTML = `
                    <div class="flex justify-between">
                        <div class="text-sm">
                            <p class="font-medium text-gray-900">${topic.name}</p>
                        </div>
                        <div class="text-sm text-gray-500">
                            <p>${topic.messageCount} messages</p>
                        </div>
                    </div>
                `;
                
                elements.topicsList.appendChild(topicItem);
            });
        }
    } catch (error) {
        console.error('Error fetching topics:', error);
    }
}

// Handle sending a message
async function handleSendMessage(event) {
    event.preventDefault();
    
    const formData = new FormData(elements.messageForm);
    const topic = formData.get('topic');
    const key = formData.get('key');
    const value = formData.get('value');
    
    try {
        elements.messageForm.querySelector('button[type="submit"]').disabled = true;
        
        // Prepare the message data
        const messageData = {
            topic: topic,
            key: key,
            value: value,
            headers: {}
        };
        
        // Get API key from localStorage or prompt user
        let apiKey = localStorage.getItem('kafkaGatewayApiKey');
        if (!apiKey) {
            apiKey = prompt('Please enter your API key:');
            if (apiKey) {
                localStorage.setItem('kafkaGatewayApiKey', apiKey);
            }
        }
        
        // Make the API call to send the message
        const response = await fetch(config.sendMessageEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${apiKey}`
            },
            body: JSON.stringify(messageData)
        });
        
        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.error || `HTTP error! status: ${response.status}`);
        }
        
        // Show success message
        alert(`Message sent successfully to topic: ${topic}`);
        
        // Reset form
        elements.messageForm.reset();
        
    } catch (error) {
        console.error('Error sending message:', error);
        alert('Error sending message: ' + error.message);
    } finally {
        elements.messageForm.querySelector('button[type="submit"]').disabled = false;
    }
}

// Initialize dashboard when page loads
document.addEventListener('DOMContentLoaded', async () => {
    await initDashboard();
});