
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM gcr.io/google-containers/debian-base-amd64:0.3
ENV DEBIAN_FRONTEND noninteractive
RUN clean-install ca-certificates
ADD cluster-autoscaler cluster-autoscaler
ADD run.sh run.sh
CMD ./run.sh
