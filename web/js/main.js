const { createApp, ref, onMounted } = Vue;
const { ElMessage } = ElementPlus;

const App = {
    setup() {
        const dialogVisible = ref(false);
        const input = ref('');
        const message = ref('');
        const currentSentence = ref({});
        const optionsAnalysis = ref([
            { value: 'literary', label: '文学表达分析' },
            { value: 'logic', label: '逻辑与含义拆解' },
            { value: 'emotion', label: '情绪与语气判断' },
            { value: 'context', label: '口语/沟通场景解读' },
            { value: 'learning', label: '学习与思考角度延展' },
        ]);
        const currentModelAnalysis = ref('literary');
        const configDialogVisible = ref(false);
        // 每日一句
        const getDailySentence = async () => {
            try {
                const res = await fetch('/api/single_sentence')
                const data = await res.json()
                if (data.code == 200) {
                    let obj = data.data || {}
                    currentSentence.value = obj
                    message.value = obj.content || 'No sentence found.'
                    analysisResult.value = ''
                }
            } catch (error) {}
        }
        // 添加 句子
        const handleAdd = async () => {
            if (!input.value) {
                ElMessage({
                    message: '请输入内容!',
                    type: 'warning',
                })
                return
            }
            try {
                const res = await fetch('/api/add_sentence', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        content: input.value,
                    }),
                })
                const data = await res.json()
                if (data.code == 200) {
                    let obj = data.data || {}
                    currentSentence.value = obj
                    message.value = input.value;
                    input.value = ''
                    dialogVisible.value = false;
                    configDialogVisible.value = false;
                    analysisResult.value = ''
                    getSentenceList();
                    ElMessage.success('添加成功!')
                } else {
                    ElMessage.error('添加失败!')
                }
            } catch (error) {
                console.error('Error adding sentence:', error)
            }
        }
        // 删除 句子
        const handleDelSentence = async (row) => {
            console.log('Delete Sentence Data:', row)
        
            const res = await fetch(`/api/del_sentence/${row.ID}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    id: row.ID,
                }),
            })
            const data = await res.json()
            if (data.code == 200) {
                ElMessage.success('删除成功!')
                getSentenceList();
            } else {
                ElMessage.error('添加失败!')
            }
        }
        
        // 句子列表
        const tableData = ref([]);
        const formPagination = ref({
            page: 1,
            pageSize: 10,
            total: 0,
        });
        const getSentenceList = async () => {
            try {
                let url = `/api/sentence_list?page=${formPagination.value.page}&pageSize=${formPagination.value.pageSize}`
                const res = await fetch(url, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                })
                const data = await res.json()
                if (data.code == 200) {
                    let list = data.data.list || [];
                    tableData.value = list
                    console.log('Sentence List:', tableData.value)
                    formPagination.value.total = data.data.total || 0
                    formPagination.value.page = data.data.page || 1
                    formPagination.value.pageSize = data.data.pageSize || 10
                }
            } catch (error) {
                console.error('Error fetching sentence list:', error)
            }
        }
        const handlePagination = (value) => {
            formPagination.value.pageNum = value || 1;
            getSentenceList()
        }
    
        // 模型列表（可用）
        const options = ref([])
        const currentModel = ref('')
        const getOllamaModelsList = async () => {
            try {
                const res = await fetch('/api/ollama_models')
                const data = await res.json()
                if (data.code == 200) {
                    let models = data.data.models || [];
                    currentModel.value = models.length > 0 ? models[0].name : ''
                    options.value = models.map(model => ({
                        value: model.name,
                        label: model.name,
                    }))
                }
                console.log('Ollama Models List:', data)
            } catch (error) {
                console.error('Error fetching Ollama models:', error)
            }
        }
        // 删除模型
        const deleteModelDialogVisible = ref(false)
        const currentModelToDelete = ref('')
        const handleDelModel = async () => {
            if (!currentModelToDelete.value) {
                ElMessage({
                    message: '请选择要删除的模型!',
                    type: 'warning',
                })
                return
            }
            try {
                const res = await fetch('/api/ollama_delete_model', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        modelName: currentModelToDelete.value,
                    }),
                })
                const data = await res.json()
                if (data.code == 200) {
                    ElMessage.success('模型删除成功!')
                    deleteModelDialogVisible.value = false
                    getOllamaModelsList()
                } else {
                    ElMessage.error('模型删除失败!')
                }
            } catch (error) {
                console.error('Error deleting Ollama model:', error)
            }
        }
        // 拉取模型
        const pullModelDialogVisible = ref(false)
        const pullModelName = ref('')
        const pullLoading = ref(false)
        const pullError = ref('')
        const handlePullModel = async () => {
            if (!pullModelName.value) {
                ElMessage({
                    message: '请输入要拉取的模型名称!',
                    type: 'warning',
                })
                return
            }
            try {
                const source = new EventSource(`/api/pull/stream?modelName=${pullModelName.value}`)
                pullError.value = ''
                pullLoading.value = true
                source.onmessage = (e) => {
                    const data = JSON.parse(e.data)
                    if (data.status === "success") {
                        pullLoading.value = false
                        pullModelDialogVisible.value = false
                        source.close()
                        ElMessage.success('模型拉取成功!')
                        setTimeout(() => {
                            getOllamaModelsList()
                        }, 1000)
                    }
                }
                source.onerror = () => {
                    pullLoading.value = false
                    pullError.value = 'Error pulling Ollama model.'
                    source.close()
                }
            } catch (error) {
                pullError.value = 'Error pulling Ollama model.'
            }
        }
        // AI分析
        const analysisResult = ref('')
        const loading = ref(false)
        const getOllamaAnalysis = async () => {
            if (!currentModel.value) {
                ElMessage({
                    message: '请选择模型!',
                    type: 'warning',
                })
                return
            }
            if (!message.value) {
                alert('Please fetch a sentence first.')
                return
            }
            try {
                loading.value = true
                analysisResult.value = ''
                const res = await fetch('/api/ollama_generate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        model: currentModel.value,
                        sentence: message.value,
                        analysis_type: currentModelAnalysis.value,
                        sentence_id: currentSentence.value.ID || 0,
                    }),
                })
                const data = await res.json()
                if (data.code == 200) {
                    analysisResult.value = data.data.response || ''
                }
                loading.value = false;
            } catch (error) {
                loading.value = false
                console.error('Error fetching Ollama analysis:', error)
            }
        }
        
        // 导出
        const exportSentences = () => {
            fetch('/export/sentences')
                .then(res => res.blob())
                .then(blob => {
                    const url = URL.createObjectURL(blob)
                    const a = document.createElement('a')
                    a.href = url
                    a.download = 'sentences.xlsx'
                    a.click()
                    URL.revokeObjectURL(url)
                })
        }
        
        onMounted(() => {
            getDailySentence()
            getOllamaModelsList()
            getSentenceList()
        })
        return {
            message,
            input,
            options,
            currentModel,
            analysisResult,
            dialogVisible,
            optionsAnalysis,
            currentModelAnalysis,
            loading,
            deleteModelDialogVisible,
            currentModelToDelete,
            pullModelDialogVisible,
            pullModelName,
            pullLoading,
            pullError,
            configDialogVisible,
            tableData,
            formPagination,
            handlePullModel,
            handleDelModel,
            handleAdd,
            getDailySentence,
            getOllamaAnalysis,
            exportSentences,
            handleDelSentence,
            handlePagination,
        }
    }
}

createApp(App).use(ElementPlus).mount('#app');