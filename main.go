// Copyright (c) 2026 SpaceTent
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    
    "github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var VERSION = "Development"
var BUILDTIME = ""

func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "k8sec: TUI-based Kubernetes Secret editor\n")
        fmt.Fprintf(os.Stderr, "Version: %s\n", VERSION)
        fmt.Fprintf(os.Stderr, "BuildTime: %s\n", BUILDTIME)
        fmt.Fprintf(os.Stderr, "\nUsage of k8sec:\n")
        flag.PrintDefaults()
    }
    
    nsFlag := flag.String("n", "", "namespace to use (overrides NAMESPACE env var, defaults to 'default')")
    flag.Parse()
    
    client := mustClient()
    ns := getNamespace(*nsFlag)
    secrets := listSecrets(client, ns)
    selected, add := selectSecret(secrets)
    
    if add {
        addSecret(client, ns)
    } else {
        editSecret(client, ns, selected)
    }
}

func mustClient() *kubernetes.Clientset {
    kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    must(err)
    client, err := kubernetes.NewForConfig(config)
    must(err)
    return client
}

func getNamespace(flagNS string) string {
    if flagNS != "" {
        return flagNS
    }
    ns := os.Getenv("NAMESPACE")
    if ns == "" {
        ns = "default"
    }
    return ns
}

func listSecrets(client *kubernetes.Clientset, ns string) []v1.Secret {
    list, err := client.CoreV1().Secrets(ns).List(context.TODO(), metav1.ListOptions{})
    must(err)
    
    return list.Items
}

func selectSecret(secrets []v1.Secret) (v1.Secret, bool) {
    if len(secrets) == 0 {
        fmt.Println("No secrets found in namespace.")
        fmt.Println("Press 'a' to add a new secret, or 'ctrl+c' to quit.")
    }
    
    items := make([]list.Item, len(secrets))
    for i, s := range secrets {
        items[i] = item{secret: s}
    }
    
    m := model{
        list: list.New(items, list.NewDefaultDelegate(), 0, 0),
    }
    m.list.Title = "Select a secret to edit ('a' to add new)"
    
    p := tea.NewProgram(m, tea.WithAltScreen())
    finalModel, err := p.Run()
    must(err)
    
    res := finalModel.(model)
    if res.addMode {
        return v1.Secret{}, true
    }
    
    sel := res.selected
    if sel == nil {
        fmt.Println("No secret selected.")
        os.Exit(0)
    }
    
    return *sel, false
}

func editSecret(client *kubernetes.Clientset, ns string, secret v1.Secret) {
    env := secretToEnv(secret)
    
    tmpFile := writeTemp(env)
    defer os.Remove(tmpFile)
    
    openEditor(tmpFile)
    
    updatedEnv := readEnv(tmpFile)
    
    if !hasChanges(secret.Data, updatedEnv) {
        fmt.Println("no changes")
        return
    }
    
    updateSecret(client, ns, secret, updatedEnv)
}

func hasChanges(old, new map[string][]byte) bool {
    if len(old) != len(new) {
        return true
    }
    for k, v := range old {
        newV, ok := new[k]
        if !ok || string(v) != string(newV) {
            return true
        }
    }
    return false
}

func secretToEnv(secret v1.Secret) string {
    var b strings.Builder
    
    for k, v := range secret.Data {
        b.WriteString(fmt.Sprintf("%s=%s\n", k, string(v)))
    }
    
    return b.String()
}

func writeTemp(content string) string {
    tmp, err := ioutil.TempFile("", "ksec-*.env")
    must(err)
    
    err = os.Chmod(tmp.Name(), 0600)
    must(err)
    
    _, err = tmp.WriteString(content)
    must(err)
    
    tmp.Close()
    return tmp.Name()
}

func openEditor(file string) {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "nano"
    }
    
    cmd := exec.Command(editor, file)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    err := cmd.Run()
    must(err)
}

func readEnv(file string) map[string][]byte {
    out := make(map[string][]byte)
    
    f, err := os.Open(file)
    must(err)
    defer f.Close()
    
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        
        if strings.TrimSpace(line) == "" {
            continue
        }
        
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        
        key := parts[0]
        val := parts[1]
        
        out[key] = []byte(val)
    }
    
    return out
}

func updateSecret(client *kubernetes.Clientset, ns string, secret v1.Secret, data map[string][]byte) {
    secret.Data = data
    
    _, err := client.CoreV1().Secrets(ns).Update(context.TODO(), &secret, metav1.UpdateOptions{})
    must(err)
    
    fmt.Printf("Secret %s updated.\n", secret.Name)
}

func addSecret(client *kubernetes.Clientset, ns string) {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Print("Enter secret name: must be lowercase alphanumeric, '-' or '_': ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    
    if name == "" {
        fmt.Println("Secret name cannot be empty.")
        return
    }
    
    fmt.Printf("Enter namespace [%s]: ", ns)
    newNs, _ := reader.ReadString('\n')
    newNs = strings.TrimSpace(newNs)
    if newNs == "" {
        newNs = ns
    }
    
    secret := v1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: newNs,
        },
        Type: v1.SecretTypeOpaque,
        Data: make(map[string][]byte),
    }
    
    tmpFile := writeTemp("")
    defer os.Remove(tmpFile)
    
    openEditor(tmpFile)
    
    data := readEnv(tmpFile)
    if len(data) == 0 {
        fmt.Println("no changes")
        return
    }
    
    secret.Data = data
    _, err := client.CoreV1().Secrets(newNs).Create(context.TODO(), &secret, metav1.CreateOptions{})
    must(err)
    
    fmt.Printf("Secret %s created in namespace %s.\n", name, newNs)
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
