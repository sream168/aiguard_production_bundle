import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ReviewForm from '../components/ReviewForm.vue'

describe('ReviewForm', () => {
  const defaultProps = {
    form: {
      mrUrl: '',
      repoUrl: '',
      localRepoPath: '',
      sourceBranch: '',
      targetBranch: '',
      configPath: 'config.yaml',
      workspaceDir: './workspace'
    },
    availableBranches: [],
    running: false,
    pullingCode: false,
    clearingCache: false,
    showLogs: false,
    disabled: false,
    repoSuggestionMessage: '',
    repoSuggestionCandidates: [],
    repoSuggestionResolvedByMr: false
  }

  it('renders form fields', () => {
    const wrapper = mount(ReviewForm, { props: defaultProps })
    expect(wrapper.find('input[placeholder*="merge_requests"]').exists()).toBe(true)
  })

  it('emits pull-code event', async () => {
    const wrapper = mount(ReviewForm, { props: defaultProps })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('pull-code')).toBeTruthy()
  })

  it('shows repo suggestion pill when available', () => {
    const wrapper = mount(ReviewForm, {
      props: {
        ...defaultProps,
        repoSuggestionMessage: '已自动推导',
        repoSuggestionResolvedByMr: true,
        repoSuggestionCandidates: ['git@example.com:test/repo.git']
      }
    })
    expect(wrapper.text()).toContain('已自动推导')
  })
})
